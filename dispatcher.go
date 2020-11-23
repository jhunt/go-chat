package chat

import (
	"regexp"
)

type rule struct {
	pattern *regexp.Regexp
	handler Handler
}

type Dispatcher struct {
	every Handler
	on    []rule
}

func (d *Dispatcher) Every(fn Handler) {
	d.every = fn
}

func (d *Dispatcher) On(regex string, fn Handler) {
	d.on = append(d.on, rule{
		pattern: regexp.MustCompile(regex),
		handler: fn,
	})
}

func (d Dispatcher) Dispatch(msg Message) {
	if d.every != nil {
		if d.every(msg) == Handled {
			return
		}
	}

	if !msg.Addressed {
		return
	}

	for _, rule := range d.on {
		if args := rule.pattern.FindStringSubmatch(msg.Text); args != nil {
			if rule.handler(msg, args...) == Handled {
				return
			}
		}
	}
}
