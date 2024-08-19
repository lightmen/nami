package ketama

type Option func(k *Ketama)

func Replicas(r int) Option {
	return func(k *Ketama) {
		k.replicas = r
	}
}

func Hash(fn HashFunc) Option {
	return func(k *Ketama) {
		k.hash = fn
	}
}
