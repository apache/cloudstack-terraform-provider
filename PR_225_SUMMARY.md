# PR #225 Summary: CI/Bump Simulator Wait

## Overview
This PR increases the CloudStack simulator readiness wait time from 10 minutes to 20 minutes to stabilize the acceptance test matrix.

## Problem
- Acceptance test matrix experiencing frequent timeouts around 20-27 minute mark
- Jobs failing out of the "Run acceptance test" step
- Issue is consistent with CloudStack simulator readiness rather than provider logic

## Solution
- Double the simulator readiness wait from 10 minutes (20 × 30s) to 20 minutes (40 × 30s)
- Change in `.github/actions/setup-cloudstack/action.yml` line 46:
  - **Before:** `until [ $T -gt 20 ] || curl -sfL http://localhost:8080 --output /dev/null`
  - **After:** `until [ $T -gt 40 ] || curl -sfL http://localhost:8080 --output /dev/null`

## Changes Made
- **File:** `.github/actions/setup-cloudstack/action.yml`
- **Line 46:** Increased timeout from 20 to 40 iterations (10m → 20m)
- **Scope:** CI only - no provider logic changes
- **Risk:** None (CI-only change)

## Repository Setup
- **Location:** `~/Downloads/cloudstack-terraform-provider/`
- **Branch:** `pr-225` (checked out from PR #225)
- **Build Status:** ✅ Successfully built
- **Binary:** `~/go/bin/terraform-provider-cloudstack`

## Dependencies Verified
- ✅ Go 1.25.1 (required: 1.20+)
- ✅ Terraform 1.10.5 (required: 1.0.x)
- ✅ Docker 28.4.0 (for CloudStack simulator)

## Testing
To test the changes:
1. Run acceptance tests: `make testacc`
2. Monitor CI logs for the increased wait time
3. Verify tests complete within the extended timeout

## Related Issues
- Relates to #218 (stabilizes acceptance matrix for that feature PR)
- Addresses frequent CI timeouts in acceptance test matrix
