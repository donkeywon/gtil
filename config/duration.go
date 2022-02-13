package config

import "time"

type Duration time.Duration

func (d *Duration) UnmarshalYAML(unmarshal func(v interface{}) error) error {
	t := ""
	err := unmarshal(&t)
	if err != nil {
		return err
	}

	tmpD, err := time.ParseDuration(t)
	if err != nil {
		return err
	}

	*d = Duration(tmpD)

	return nil
}

func (d Duration) ToDuration() time.Duration {
	return time.Duration(d)
}
