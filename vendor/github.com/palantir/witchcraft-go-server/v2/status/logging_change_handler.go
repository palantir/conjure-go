// Copyright (c) 2020 Palantir Technologies. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package status

import (
	"context"

	"github.com/palantir/witchcraft-go-health/conjure/witchcraft/api/health"
	"github.com/palantir/witchcraft-go-health/status"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
)

func loggingHealthStatusChangeHandler() HealthStatusChangeHandler {
	return healthStatusChangeHandlerFn(func(ctx context.Context, prevStatus, currStatus health.HealthStatus) {
		logIfHealthChanged(ctx, prevStatus, currStatus)
	})
}

func logIfHealthChanged(ctx context.Context, previousHealth, newHealth health.HealthStatus) {
	previousCode, newCode := status.HealthStatusCode(previousHealth), status.HealthStatusCode(newHealth)
	if previousCode != newCode {
		params := map[string]interface{}{
			"previousHealthStatusCode": previousCode,
			"newHealthStatusCode":      newCode,
			"newHealthStatus":          newHealth.Checks,
		}
		switch newCode {
		case 200: // HEALTHY
			svc1log.FromContext(ctx).Info("Health status code changed.", svc1log.SafeParams(params))
		case 518, 520: // DEFERRING, REPAIRING
			svc1log.FromContext(ctx).Warn("Health status code changed.", svc1log.SafeParams(params))
		default:
			svc1log.FromContext(ctx).Error("Health status code changed.", svc1log.SafeParams(params))
		}
		return
	}
	svc1log.FromContext(ctx).Info("Health checks content changed without status change.", svc1log.SafeParams(map[string]interface{}{
		"statusCode":      newCode,
		"newHealthStatus": newHealth.Checks,
	}))
}
