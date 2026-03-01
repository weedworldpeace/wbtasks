package service

import (
	"app/internal/models"
	"strings"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) Cut(ent models.Entity) ([]string, error) {
	return cutStrings(ent.Data, ent.Args)
}

func cutStrings(data []string, args models.Arguments) ([]string, error) {
	res := make([]string, 0, len(data))
	for i := range data {
		splitted := strings.Split(data[i], args.D)
		if (args.S && len(splitted) > 1) || (!args.S) {
			toJoin := make([]string, 0, len(args.F))
			for _, v := range args.F {
				if len(splitted) < v {
					break
				}
				toJoin = append(toJoin, splitted[v-1])
			}
			res = append(res, strings.Join(toJoin, args.D))
		}
	}
	return res, nil
}
