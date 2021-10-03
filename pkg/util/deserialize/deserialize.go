package deserialize

import "gopkg.in/yaml.v2"

// ParseYamlConfig 通用反序列化yaml方法
func ParseYamlConfig(b []byte, object interface{}) (interface{}, error) {
	if err := yaml.Unmarshal(b, &object); err != nil {
		return nil, err
	}

	return object, nil
}
