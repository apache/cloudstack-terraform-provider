<!--
  Licensed to the Apache Software Foundation (ASF) under one
  or more contributor license agreements.  See the NOTICE file
  distributed with this work for additional information
  regarding copyright ownership.  The ASF licenses this file
  to you under the Apache License, Version 2.0 (the
  "License"); you may not use this file except in compliance
  with the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing,
  software distributed under the License is distributed on an
  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
  KIND, either express or implied.  See the License for the
  specific language governing permissions and limitations
  under the License.
-->

# Apache CloudStack Terraform Provider Security Threat Model — delta (draft)

> **Delta document.** Inherits §3, §4 B1, and §7 from
> `cloudstack-threat-model-draft.md`. Read the main model first.

## §1 Header

- **Project:** Apache CloudStack Terraform Provider
  (`apache/cloudstack-terraform-provider`) — Terraform provider that
  drives the CloudStack JSON API as a Hashicorp Terraform plugin.
- **Commit:** `f19bffc` (HEAD of `main` at draft time).
- **Date:** 2026-05-29.
- **Authors:** ASF Security team draft.
- **Status:** Draft delta over `cloudstack-threat-model-draft.md`.
- **Version binding:** as of the commit above. Provider released at
  `v0.5.x` series *(documented: `README.md` example)*.
- **Reporting:** as in the main model.
- **Provenance legend:** as in the main model.
- **Draft confidence:** 8 documented / 0 maintainer / 10 inferred.

**About the project.** A Terraform provider plugin written in Go,
built on top of `apache/cloudstack-go`. Each Terraform resource
(`cloudstack_instance`, `cloudstack_network`, `cloudstack_template`,
…) and data source maps to one or more CloudStack API calls. Runs
inside the Terraform plugin host as a separate process invoked by
the `terraform` binary over the Terraform plugin RPC protocol
*(documented: `README.md`, Terraform plugin architecture)*.

## §2 Scope and intended use

**Primary intended use.** *(documented — README)* Operators describe
their CloudStack-managed infrastructure as Terraform resources, then
`terraform apply` drives the API to converge. Used for IaC-style
infrastructure provisioning.

**Deployment shape.** A short-lived Go process launched by the
Terraform binary on the operator's workstation or a CI runner. No
long-running daemon, no listener.

**Caller expectations.** The caller (the operator / CI) is trusted to:

- supply provider config (`api_url`, `api_key`, `secret_key`,
  `http_get_only`, `timeout`) via either the provider block or env
  vars,
- not source provider config from end-user input,
- store the Terraform state file (which may contain resource IDs and
  some sensitive metadata) per Terraform's own best practices.

**Component-family table.**

| Family | Representative entry | Touches outside the process? | In this delta? |
| --- | --- | --- | --- |
| Provider config / client (`cloudstack/config.go`) | `Config.NewClient()` *(documented)* | **yes — network + creds via cloudstack-go** | yes |
| Resource implementations (`cloudstack/resource_cloudstack_*.go`, ~80 resources) | CRUD lifecycles | inherited from client | yes |
| Data sources (`cloudstack/data_source_cloudstack_*.go`, ~30 data sources) | read-only lifecycles | inherited from client | yes |
| `main.go` (provider entry) | Terraform plugin handshake | Terraform RPC | yes |
| `website/`, `scripts/`, `CHANGELOG.md` | docs + release tooling | n/a | **out of model** *(§3)* |
| Tests under `cloudstack/*_test.go` | acceptance tests | network | **out of model** *(§3)* |
| `vendor/` (if vendored) | upstream Go deps | n/a | **out of model** *(§3)* |

## §3 Out of scope (explicit non-goals)

The main model's §3 applies. **Additional** out-of-scope items
specific to the Terraform provider:

1. **Terraform state file confidentiality and storage.** Where the
   operator stores Terraform state (local, S3, Terraform Cloud, …)
   and what protections are applied are operator decisions out of
   model. *(inferred — Q1)*
2. **Server-side correctness of management-server responses.** Same
   as the Go SDK delta.
3. **Terraform plugin RPC trust.** The provider trusts the `terraform`
   binary that invoked it.
4. **`vendor/`, `website/`, `scripts/`, tests.**
5. **The four sibling repos.**
6. **Terraform itself.** Bugs in HashiCorp Terraform are upstream.
7. **HCL configuration trust.** A malicious `*.tf` file is `OUT-OF-MODEL:
   trusted-input` — operator's responsibility.

## §4 Trust boundaries and data flow

The boundary is **the Terraform plugin handshake** plus the
provider-block config. Provider config (`api_url`, `api_key`,
`secret_key`, `http_get_only`) flows from the operator's environment
into `Config.NewClient()` *(documented: `cloudstack/config.go`)*. The
client then makes B1-shape API calls; see the Go SDK delta for the
signing flow.

Crucially, *(documented: `cloudstack/config.go` line ~36)*:

```go
cs := cloudstack.NewAsyncClient(c.APIURL, c.APIKey, c.SecretKey, false)
cs.HTTPGETOnly = c.HTTPGETOnly
```

The fourth `NewAsyncClient` argument is **hardcoded to `false`** —
meaning **`InsecureSkipVerify: true` is unconditionally set** in the
TLS config of every Terraform-provider client *(documented: see the Go
SDK `cloudstack.go` line ~216)*. That is the single most security-
relevant fact in this delta and is the centerpiece of §14 Q2 below.
The provider does **not** appear to expose a `verify_ssl` provider
attribute to the operator *(inferred — needs maintainer confirmation
via Q2)*.

## §5 Assumptions about the environment

- **Host**: Linux/macOS/Windows operator workstation or CI runner;
  Terraform 1.0+ *(documented: `README.md`)*.
- **Go**: 1.20+ *(documented: `README.md`)*.
- **Filesystem**: Terraform state file under operator-chosen backend.
- **Network**: outbound HTTPS to `api_url`.
- **What the provider does not do**: no listener, no daemon, no
  global state mutation outside Terraform plugin handshake
  *(inferred — Q3)*.

## §5a Build-time and configuration variants

| Knob | Default | Stance | Effect |
| --- | --- | --- | --- |
| `api_url` | none | operator config | endpoint |
| `api_key`, `secret_key` | none | operator config | credentials |
| `http_get_only` | per provider block | when `true`, signatures land in URL — see Go SDK delta Q3 | transport choice |
| `timeout` | provider default *(inferred — Q4)* | bounds async polling | |
| TLS verification | **hardcoded off (`verifyssl=false`)** *(documented: `cloudstack/config.go` line ~36)* | **maintainer ruling required** *(inferred — Q2)* | **all HTTPS calls skip cert verification** |

## §6 Assumptions about inputs

| Entry point | Parameter | Attacker-controllable in the model? | Caller must enforce |
| --- | --- | --- | --- |
| Provider block | `api_url`, `api_key`, `secret_key`, `http_get_only`, `timeout` | **no** — operator / CI config | not from end-user input |
| HCL resource definitions | resource attributes | **no** — IaC author | not from end-user input |
| HTTP response | JSON body | trusted (B1) | typed-decoded |

## §7 Adversary model

Main-model §7 applies. **Adjustments specific to this provider**:

- "Unauthenticated network peer reaching `:8080`" is upstream.
- **A network adversary between the operator / CI runner and the
  CloudStack management server can man-in-the-middle the API
  conversation undetected**, because TLS verification is disabled at
  the provider layer. Whether this is a `VALID` report (provider
  should change) or `BY-DESIGN: property-disclaimed` (provider
  intentionally delegates TLS to the operator's network) is Q2.

## §8 Security properties the provider provides

### T1 — HMAC-SHA1 signature via `cloudstack-go`

- **Property.** Each API call carries a valid HMAC-SHA1 signature per
  main-model §8 P1, via the Go SDK.
- **Conditions / violation / severity.** As main model §8 P1.

### T2 — Terraform plugin handshake

- **Property.** The provider speaks only the Terraform plugin RPC
  protocol on stdin/stdout to its parent `terraform` process.
- **Conditions.** Invoked by Terraform.
- **Violation symptom.** Provider accepts commands from a non-Terraform
  caller.
- **Severity.** Security-critical, `VALID`.

## §9 Security properties the provider does *not* provide

- **No TLS verification of the management-server cert** — TLS-verify
  is hardcoded off in `cloudstack/config.go` *(documented)*. This is
  the headline non-property and the single biggest §14 question.
- **No protection of `secret_key` in Terraform state.** Provider-level
  secrets may end up in plan output and state.
- **No re-validation of management-server response correctness.**
- **No defence against malicious Terraform HCL.** A `*.tf` file with
  destructive resource definitions executes as written.

### False-friend properties

- **The provider running over HTTPS is not "TLS-protected" in the
  authentication sense** because cert verification is off. Any
  on-path attacker can MitM, present any cert, see all signed
  requests, and replay them within their `expires` window.

## §10 Downstream responsibilities

The operator / CI MUST:

1. Run the provider only over a network where the
   operator-to-management-server path is already trusted (private
   network, VPN, internal CI runner inside the CloudStack control
   plane), because the provider does not verify the management-server
   cert. *(pending Q2 — if the PMC adds a `verify_ssl` knob, this
   responsibility narrows.)*
2. Treat the Terraform state file as sensitive — `secret_key` and
   resource secrets may live there.
3. Source `api_key` / `secret_key` from a secret manager, not from
   the HCL file.
4. Match provider version to the management-server version.
5. Apply `http_get_only` only when URL logs cannot leak the signature.

## §11 Known misuse patterns

- Running the provider over an untrusted network (public Internet)
  expecting TLS verification.
- Hardcoding `secret_key` in HCL or env vars committed to a repo.
- Storing Terraform state in a backend without per-team ACLs.
- Sharing one set of `api_key` / `secret_key` across multiple CI
  jobs without isolation.

## §11a Known non-findings (recurring false positives)

- **"HMAC-SHA1 — SHA1 is deprecated."** → `KNOWN-NON-FINDING` per main
  model §11a.
- **"`InsecureSkipVerify: true` is hardcoded in `Config.NewClient`."**
  This is the *actual provider behaviour* *(documented:
  `cloudstack/config.go` line ~36)*. Whether it is a `VALID` finding
  or `BY-DESIGN: property-disclaimed` is Q2. **Do not silently
  suppress.**
- **"Secrets appear in Terraform state."** Out of model — Terraform
  state security is operator's job. → `OUT-OF-MODEL: trusted-input`.
- **"Tests in `cloudstack/*_test.go` have weak input handling."**
  Out of model. → `OUT-OF-MODEL: unsupported-component`.

## §12 Conditions that would change this delta

- A `verify_ssl` (or `insecure_skip_verify`) provider attribute is
  added — would change Q2.
- A switch in default to verifying TLS — would change §9 first
  bullet.
- A new resource that introduces a new credential-bearing state shape
  (e.g. a `cloudstack_user_credentials` resource that persists secrets
  outside Terraform state).
- Change in signing algorithm at the main-model layer.

## §13 Triage dispositions

Use the same table as the main model.

## §14 Open questions for the maintainers

**Q1.** Out of scope: Terraform state file storage and encryption.
Confirm.

**Q2.** **Highest-leverage question in this delta.** The provider
calls `cloudstack.NewAsyncClient(c.APIURL, c.APIKey, c.SecretKey,
false)` in `cloudstack/config.go` line ~36 — hardcoding
`verifyssl=false` and thus `InsecureSkipVerify: true` on every HTTP
call.

- (a) Is this intentional (the provider assumes operators run only
  on a trusted control-plane network and explicitly opts out of TLS
  verification), and should the §9 first bullet stay?
- (b) Or should the provider expose a `verify_ssl` attribute and
  default it to `true`, in which case the current hardcoding is a
  `VALID` finding to fix?

The text of §9 and §10 assumes (a) — please confirm or correct.
*(maps to §5a, §9, §10, §11a)*

**Q3.** Confirm §5 negative side-effect inventory: no env-var
consumption beyond Terraform plugin standard, no signal handlers, no
global mutation.

**Q4.** Confirm `timeout` default for the cloudstack-go async client
as used by the provider.

**Q5.** Where in the provider's docs is the credential-handling
responsibility documented for HCL authors?

**Q6.** Does the provider redact `secret_key` from Terraform plan
output? From `terraform show`?

**Q7.** Does the provider integrate with a Vault or other
secret-manager pattern, or is the operator expected to wire that in
externally?

**Q8.** Meta — should this delta live at `docs/threat-model.md` in
`apache/cloudstack-terraform-provider`, or in the website tree?

**Q9.** When the main model's signing algorithm changes (e.g. v3 →
v4, SHA1 → SHA256), what is the release-cadence commitment to update
the provider's vendored `cloudstack-go`?

**Q10.** Confirm the unsupported-component list (tests, website,
scripts, vendor).
