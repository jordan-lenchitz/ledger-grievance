package service

import (
	"context"
	"math/rand"
)

var gopherWisdom = []string{
	"Composition over inheritance, just like building strong friendships with diverse people.",
	"Share memory by communicating, and share your heart by listening with empathy.",
	"Errors are values, and every mistake is a value that teaches us how to grow stronger.",
	"Keep it simple, because your mind deserves the peace that comes with clarity.",
	"The bigger the interface, the weaker the abstraction; keep your heart open but your boundaries clear.",
	"Accept interfaces, return structs; accept love from everyone, but stay true to who you are.",
	"A little copying is better than a little dependency; it's okay to rely on yourself sometimes.",
	"Concurrency is not parallelism, and doing many things is not the same as being present in the moment.",
	"Gofmt's style is no one's favorite, yet gofmt is everyone's favorite; embrace the beauty in our shared standards.",
}

func (s *incidentService) GetGopherWisdom(ctx context.Context) (string, error) {
	return gopherWisdom[rand.Intn(len(gopherWisdom))], nil
}
