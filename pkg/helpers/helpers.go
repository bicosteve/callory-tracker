package helpers

import "path/filepath"

func LoadTemplate(template string) (string, error) {
	temp, err := filepath.Abs(template)
	if err != nil {
		return "", err
	}

	return temp, nil
}

func LoadEnv(name string) (string, error) {
	file, err := filepath.Abs(name)
	if err != nil {
		return "", err
	}

	return file, nil
}

//func LoadTemplates(dir string) (string, error) {
//	path, err := filepath.Abs(dir)
//	if err != nil {
//		return "", err
//	}
//
//	return path, nil
//}
