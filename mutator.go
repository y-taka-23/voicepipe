package main

func replaceEnv(df Dockerfile, key, value string) *Dockerfile {
	sts := []Statement{}
	for _, st := range df.Statements {
		if env, ok := st.(*Env); ok {
			vs := map[string]string{}
			for k, v := range env.Variables {
				if k == key {
					vs[key] = value
				} else {
					vs[k] = v
				}
			}
			sts = append(sts, Env{Variables: vs})
		} else {
			sts = append(sts, st)
		}
	}
	return &Dockerfile{Statements: sts}
}
