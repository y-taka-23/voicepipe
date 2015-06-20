package main

func replaceEnv(df Dockerfile, k, v string) *Dockerfile {
	sts := []Statement{}
	for _, st := range df.Statements {
		if x, ok := st.(*Env); ok {
			x.Variables[k] = v
			sts = append(sts, x)
		} else {
			sts = append(sts, st)
		}
	}
	return &Dockerfile{Statements: sts}
}
