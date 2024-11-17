package scheduler

import (
	"log/slog"
	"strings"

	"github.com/adalbertjnr/downscaler-operator/api/v1alpha1"
	"github.com/adalbertjnr/downscaler-operator/internal/pkgerrors"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

const (
	Spec = "spec"

	Schedule = "schedule"
	TimeZone = "timeZone"

	DownscalerOptions = "downscalerOptions"
	TimeRules         = "timeRules"
	Rules             = "rules"
	Namespaces        = "namespaces"
	UpscaleTime       = "upscaleTime"
	DownscaleTime     = "downscaleTime"
)

func (s *Downscaler) Validate() bool {
	valid := true

	var validationErrors []error

	processScheduleFields(&s.app.Spec.Schedule, &validationErrors)
	processDownscalerOptions(&s.app.Spec.DownscalerOptions, &validationErrors)

	if len(validationErrors) > 0 {
		for _, err := range validationErrors {
			slog.Error("validation failed", "err", err)
		}
		return !valid
	}

	return valid
}

func processScheduleFields(schedule *v1alpha1.Schedule, validationErrors *[]error) {
	if schedule == nil {
		err := field.Invalid(field.NewPath(Spec).Child(Schedule), schedule, pkgerrors.ErrNilInclude.Error())
		*validationErrors = append(*validationErrors, err)
		return
	}
	if schedule.TimeZone == "" || len(strings.Split(schedule.TimeZone, "/")) == 1 {
		err := field.Invalid(field.NewPath(Spec).Child(Schedule).Child(TimeZone), schedule.TimeZone, pkgerrors.ErrMalformedTimeZone.Error())
		*validationErrors = append(*validationErrors, err)
	}
}

func processDownscalerOptions(options *v1alpha1.DownscalerOptions, validationErrors *[]error) {
	childBase := field.NewPath(Spec).Child(DownscalerOptions)

	if options == nil {
		err := field.Invalid(childBase, options, pkgerrors.ErrNilInclude.Error())
		*validationErrors = append(*validationErrors, err)
		return
	}

	if options.TimeRules == nil {
		err := field.Invalid(childBase.Child(TimeRules), options.TimeRules, pkgerrors.ErrTimeRulesBlockNotProvided.Error())
		*validationErrors = append(*validationErrors, err)
		return
	}

	if options.TimeRules.Rules == nil {
		err := field.Invalid(childBase.Child(TimeRules).Child(Rules), options.TimeRules.Rules, pkgerrors.ErrRulesNotProvided.Error())
		*validationErrors = append(*validationErrors, err)
		return
	}

	for index, rule := range options.TimeRules.Rules {
		childRule := childBase.Child(TimeRules).Index(index)

		if len(rule.Namespaces) == 0 {
			err := field.Invalid(childRule.Child(Namespaces), rule.Namespaces, pkgerrors.ErrEmptyNamespaces.Error())
			*validationErrors = append(*validationErrors, err)
		}

		if len(strings.Split(rule.UpscaleTime, ":")) == 1 {
			err := field.Invalid(childRule.Child(Namespaces).Child(UpscaleTime), rule.UpscaleTime, pkgerrors.ErrMalforedUpscaleTime.Error())
			*validationErrors = append(*validationErrors, err)
		}

		if len(strings.Split(rule.DownscaleTime, ":")) == 1 {
			err := field.Invalid(childRule.Child(Namespaces).Child(DownscaleTime), rule.DownscaleTime, pkgerrors.ErrMalforedDownscaleTime.Error())
			*validationErrors = append(*validationErrors, err)
		}

	}
}
