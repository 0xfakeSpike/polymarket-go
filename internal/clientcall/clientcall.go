// Package clientcall invokes *polymarket.Client methods by name with JSON arguments (CLI/MCP bridge).
package clientcall

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/0xfakeSpike/polymarket-go/polymarket"
)

var contextType = reflect.TypeOf((*context.Context)(nil)).Elem()

// ClientPtrType is reflect.TypeOf((*polymarket.Client)(nil)) for listing/help.
func ClientPtrType() reflect.Type { return reflect.TypeOf((*polymarket.Client)(nil)) }

// ListClientMethods returns sorted exported method names on *polymarket.Client.
func ListClientMethods() []string { return ListExportedMethods(ClientPtrType()) }

// ListExportedMethods returns sorted exported method names on *polymarket.Client (pointer receiver).
func ListExportedMethods(clientTyp reflect.Type) []string {
	if clientTyp.Kind() != reflect.Ptr {
		return nil
	}
	var names []string
	for i := 0; i < clientTyp.NumMethod(); i++ {
		m := clientTyp.Method(i)
		if m.PkgPath != "" {
			continue
		}
		names = append(names, m.Name)
	}
	sort.Strings(names)
	return names
}

// Invoke calls client.MethodName with JSON array args (one element per non-context parameter).
// context.Context parameters are filled with context.Background() and do not consume JSON slots.
// Function-typed parameters and non-empty interfaces (e.g. WebSocket handlers) are not supported.
func Invoke(client any, method string, argsJSON json.RawMessage) (any, error) {
	v := reflect.ValueOf(client)
	if !v.IsValid() || v.Kind() != reflect.Ptr || v.IsNil() {
		return nil, fmt.Errorf("client must be non-nil pointer")
	}
	m := v.MethodByName(method)
	if !m.IsValid() {
		return nil, fmt.Errorf("unknown method %q (see: pmctl methods)", method)
	}
	mt := m.Type()
	args, err := buildCallArgs(mt, argsJSON)
	if err != nil {
		return nil, fmt.Errorf("args for %q: %w", method, err)
	}
	out := m.Call(args)
	if len(out) == 0 {
		return nil, nil
	}
	last := out[len(out)-1]
	if last.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		if !last.IsNil() {
			return nil, last.Interface().(error)
		}
		out = out[:len(out)-1]
	}
	return encodeReturns(out), nil
}

func buildCallArgs(mt reflect.Type, argsJSON json.RawMessage) ([]reflect.Value, error) {
	var slots []json.RawMessage
	if len(strings.TrimSpace(string(argsJSON))) == 0 || string(argsJSON) == "null" {
		slots = nil
	} else {
		if err := json.Unmarshal(argsJSON, &slots); err != nil {
			return nil, fmt.Errorf("args must be a JSON array (e.g. [\"id\",1]); got: %w", err)
		}
	}
	var args []reflect.Value
	idx := 0
	for i := 0; i < mt.NumIn(); i++ {
		in := mt.In(i)
		if in == contextType {
			args = append(args, reflect.ValueOf(context.Background()))
			continue
		}
		switch in.Kind() {
		case reflect.Func:
			return nil, fmt.Errorf("parameter %d (%s): function arguments are not supported via JSON; use the Go SDK", i, in)
		case reflect.Interface:
			if in.NumMethod() > 0 {
				return nil, fmt.Errorf("parameter %d (%s): interface with methods (e.g. handlers) is not supported via JSON; use the Go SDK", i, in)
			}
		}
		if idx >= len(slots) {
			return nil, fmt.Errorf("missing arg %d for type %s (need %d JSON values for non-context params)", idx, in, countNonContextParams(mt))
		}
		arg, err := decodeArg(in, slots[idx])
		if err != nil {
			return nil, fmt.Errorf("arg %d (%s): %w", idx, in, err)
		}
		args = append(args, arg)
		idx++
	}
	if idx != len(slots) {
		return nil, fmt.Errorf("too many args: got %d JSON values, need %d", len(slots), countNonContextParams(mt))
	}
	return args, nil
}

func countNonContextParams(mt reflect.Type) int {
	n := 0
	for i := 0; i < mt.NumIn(); i++ {
		if mt.In(i) != contextType {
			n++
		}
	}
	return n
}

func decodeArg(t reflect.Type, raw json.RawMessage) (reflect.Value, error) {
	if len(raw) == 0 || string(raw) == "null" {
		switch t.Kind() {
		case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice:
			return reflect.Zero(t), nil
		}
		return reflect.Value{}, fmt.Errorf("null not allowed for %s", t)
	}

	switch t.Kind() {
	case reflect.Ptr:
		if string(raw) == "null" {
			return reflect.Zero(t), nil
		}
		elem := t.Elem()
		ev := reflect.New(elem)
		if err := json.Unmarshal(raw, ev.Interface()); err != nil {
			return reflect.Value{}, err
		}
		return ev, nil
	case reflect.Interface:
		var v any
		if err := json.Unmarshal(raw, &v); err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(v), nil
	default:
		pv := reflect.New(t)
		if err := json.Unmarshal(raw, pv.Interface()); err != nil {
			return reflect.Value{}, err
		}
		return pv.Elem(), nil
	}
}

func encodeReturns(outs []reflect.Value) any {
	switch len(outs) {
	case 0:
		return nil
	case 1:
		return valueToJSONable(outs[0])
	default:
		arr := make([]any, len(outs))
		for i, o := range outs {
			arr[i] = valueToJSONable(o)
		}
		return map[string]any{"returns": arr}
	}
}

func valueToJSONable(v reflect.Value) any {
	if !v.IsValid() {
		return nil
	}
	if v.Kind() == reflect.Interface && v.IsNil() {
		return nil
	}
	if !v.CanInterface() {
		return fmt.Sprintf("%v", v)
	}
	x := v.Interface()
	if _, err := json.Marshal(x); err != nil {
		return fmt.Sprintf("%v", x)
	}
	return x
}

// MethodHelp returns a one-line signature for documentation (best-effort).
func MethodHelp(name string) (string, error) {
	m, ok := ClientPtrType().MethodByName(name)
	if !ok {
		return "", errors.New("unknown method")
	}
	var b strings.Builder
	b.WriteString(name)
	b.WriteByte('(')
	mt := m.Type
	for i := 0; i < mt.NumIn(); i++ {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(mt.In(i).String())
	}
	b.WriteString(")")
	if mt.NumOut() > 0 {
		b.WriteString(" ")
		outs := make([]string, mt.NumOut())
		for i := 0; i < mt.NumOut(); i++ {
			outs[i] = mt.Out(i).String()
		}
		b.WriteString(strings.Join(outs, ", "))
	}
	return b.String(), nil
}
