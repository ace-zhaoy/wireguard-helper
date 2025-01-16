package utils

import (
	"github.com/mitchellh/mapstructure"
	"time"
)

func SetValue[T comparable](o *T, v, defaultV T) {
	if v == *new(T) {
		*o = defaultV
	} else {
		*o = v
	}
}

func SetDefaultValue[T comparable](o *T, defaultV T) {
	if *o == *new(T) {
		*o = defaultV
	}
}

func ToStruct(input any, output any) error {
	config := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		TagName:          "json",
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			mapstructure.StringToTimeHookFunc(time.RFC3339),
			mapstructure.StringToIPHookFunc(),
			mapstructure.StringToIPNetHookFunc(),
		),
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}
