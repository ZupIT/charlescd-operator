package controllers

import "github.com/tiagoangelozup/charles-alpha/internal/predicate"

type Predicates struct {
	PredicateRepoStatus *predicate.RepoStatus
	PredicateModule     *predicate.Module
}
