package config

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func mock() string {
	return `
		[merge]
		order = ["datacenter", "environment"]

		[[profiles]]
		name        = "aws-va.production.worker"
		environment = "production"
		datacenter  = "aws-va"
		color       = "blue"

		[[profiles]]
		name        = "atlas-glete.production.worker"
		environment = "production"
		datacenter  = "atlas-glete"
		color       = "blue"

		[log]
		format 	= "human"
		size    = 1
		path    = "/var/log/application-${color}"

		[production.log]
		format = "kibanajson"

		[atlas-glete.staging.log]
		format = "human"
		size   = 0

		[[tags]]
		key   = "key 1"
		value = "value"

		[[tags]]
		key   = "key 2"
		value = ["value 1", "value 2"]

		[int-slice]
		value = [1, 2, 3]

		[bool]
		value = true

		[string-slice]
		value              = ["a", "b", "c"]
		interpolated-value = ["a-${ color }"]

		[map-string-string]
		  [[map-string-string.tags]]
		  key   = "key"
		  value = "value"

		  [[map-string-string.interpolated-tags]]
		  key   = "key"
		  value = "value-${color }"
	`
}

func TestInitialization(t *testing.T) {
	client, err := New("aws-va.production.worker", mock())
	if err != nil {
		t.Error(errors.Wrap(err, "expected to client initialization to not give error").Error())
		t.FailNow()
	}

	if !reflect.DeepEqual(client.order, []string{"datacenter", "environment"}) {
		t.Error("error invalid merge.order")
		t.FailNow()
	}

	if !reflect.DeepEqual(client.profile, map[string]string{
		"name":        "aws-va.production.worker",
		"environment": "production",
		"datacenter":  "aws-va",
		"color":       "blue",
	}) {
		t.Error("error invalid profile")
		t.FailNow()
	}
}

func TestGetString(t *testing.T) {
	client, err := New("aws-va.production.worker", mock())
	if err != nil {
		t.Error(errors.Wrap(err, "expected to client initialization to not give error").Error())
		t.FailNow()
	}

	if client.GetString("log.format") != "kibanajson" {
		t.Error("invalid log.format")
		t.FailNow()
	}

	if client.GetInt("log.size") != 1 {
		t.Error("invalid log.size")
		t.FailNow()
	}

	if client.GetString("log.path") != "/var/log/application-blue" {
		t.Error("invalid log.path")
		t.FailNow()
	}
}

func TestReplaceProfileVars(t *testing.T) {
	client, err := New("aws-va.production.worker", mock())
	if err != nil {
		t.Error(errors.Wrap(err, "expected to client initialization to not give error").Error())
		t.FailNow()
	}

	expected := "aws-va - production"
	if !reflect.DeepEqual(client.replaceProfileVars("${datacenter} - ${environment}"), expected) {
		t.Errorf(
			"got %s, expected %s",
			client.replaceProfileVars("${datacenter} - ${environment}"),
			expected,
		)
		t.Error("invalid result")
	}
}

func TestGenCombinations(t *testing.T) {
	client, err := New("aws-va.production.worker", mock())
	if err != nil {
		t.Error(errors.Wrap(err, "expected to client initialization to not give error").Error())
		t.FailNow()
	}

	expected := []string{
		"aws-va.production.log.format",
		"production.log.format",
		"aws-va.log.format",
		"log.format",
	}
	if !reflect.DeepEqual(expected, client.genCombinations("log.format")) {
		t.Error("invalid combinations")
	}
}

func TestGetRawKey(t *testing.T) {
	client, err := New("aws-va.production.worker", mock())
	if err != nil {
		t.Error(errors.Wrap(err, "expected to client initialization to not give error").Error())
		t.FailNow()
	}

	expected := []interface{}{
		map[string]interface{}{
			"key":   "key 1",
			"value": "value",
		},
		map[string]interface{}{
			"key": "key 2",
			"value": []interface{}{
				"value 1",
				"value 2",
			},
		},
	}
	if !reflect.DeepEqual(expected, client.GetRawKey("tags")) {
		t.Error("invalid result")
	}
}

func TestGetIntSlice(t *testing.T) {
	client, err := New("aws-va.production.worker", mock())
	if err != nil {
		t.Error(errors.Wrap(err, "expected to client initialization to not give error").Error())
		t.FailNow()
	}

	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(expected, client.GetIntSlice("int-slice.value")) {
		t.Error("invalid result")
	}
}

func TestGetBool(t *testing.T) {
	client, err := New("aws-va.production.worker", mock())
	if err != nil {
		t.Error(errors.Wrap(err, "expected to client initialization to not give error").Error())
		t.FailNow()
	}

	expected := true
	if !reflect.DeepEqual(expected, client.GetBool("bool.value")) {
		t.Error("invalid result")
	}
}

func TestGetStringSlice(t *testing.T) {
	client, err := New("aws-va.production.worker", mock())
	if err != nil {
		t.Error(errors.Wrap(err, "expected to client initialization to not give error").Error())
		t.FailNow()
	}

	expected := []string{"a", "b", "c"}
	if !reflect.DeepEqual(expected, client.GetStringSlice("string-slice.value")) {
		t.Error("invalid result")
	}

	expected = []string{"a-blue"}
	if !reflect.DeepEqual(expected, client.GetStringSlice("string-slice.interpolated-value")) {
		t.Error("invalid result")
	}
}

func TestGetMapStringString(t *testing.T) {
	client, err := New("aws-va.production.worker", mock())
	if err != nil {
		t.Error(errors.Wrap(err, "expected to client initialization to not give error").Error())
		t.FailNow()
	}

	expected := map[string]string{"key": "value"}
	if !reflect.DeepEqual(expected, client.GetMapStringString("map-string-string.tags")) {
		t.Error("invalid result")
	}

	expected = map[string]string{"key": "value-blue"}
	if !reflect.DeepEqual(expected, client.GetMapStringString("map-string-string.interpolated-tags")) {
		t.Error("invalid result")
	}
}
