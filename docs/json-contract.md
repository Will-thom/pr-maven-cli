# JSON Contract

PR Maven CLI emits JSON for CI systems, bots, and coding agents that need deterministic Maven failure context.

Run:

```bash
prmaven why -project . -format json
```

The top-level JSON value is a `Report`.

## Machine-Readable Schema

A JSON Schema for the current report contract is versioned at:

```text
schema/prmaven-report.schema.json
```

CI systems can validate generated reports with any JSON Schema Draft 2020-12 compatible validator. Example with `ajv-cli`:

```bash
prmaven why -project . -format json > prmaven-report.json || true
npx ajv-cli validate -s schema/prmaven-report.schema.json -d prmaven-report.json
```

## Compatibility Expectations

The JSON contract is intended to be stable for consumers that automate Maven PR triage.

Compatibility expectations:

- Existing fields should not be renamed or removed without an explicit compatibility note.
- New fields may be added in future releases.
- Consumers should ignore fields they do not understand.
- Path fields use `/` separators, including on Windows.
- `projectRoot` is an absolute local path and may differ between machines and CI workers.
- Empty optional strings such as `failureType` and `message` may be omitted.

## Report

| Field | Type | Description |
| --- | --- | --- |
| `projectRoot` | string | Absolute path to the analyzed Maven workspace. |
| `summary` | object | Aggregate counts for modules, reports, and findings. |
| `modules` | array of `Module` | Maven modules discovered from `pom.xml` files. |
| `findings` | array of `Finding` | Actionable Maven failure or quality findings. |

## Summary

| Field | Type | Description |
| --- | --- | --- |
| `moduleCount` | integer | Number of discovered Maven modules, including the root project. |
| `reportCount` | integer | Number of report or log artifacts analyzed. |
| `findingCount` | integer | Number of emitted findings. |

## Module

| Field | Type | Description |
| --- | --- | --- |
| `name` | string | Maven artifact id when available, otherwise a path-derived module name. |
| `path` | string | Module path relative to `projectRoot`; root is `"."`. |
| `pom` | string | Module `pom.xml` path relative to `projectRoot`. |

## Finding

| Field | Type | Description |
| --- | --- | --- |
| `id` | string | Stable, lowercase identifier derived from module, source/test context, and failure kind. |
| `module` | string | Human-readable module name. |
| `modulePath` | string | Module path relative to `projectRoot`; root is `"."`. |
| `reportPath` | string | Report or log artifact path relative to `projectRoot`. |
| `reportKind` | string | Source kind, such as `surefire`, `failsafe`, `checkstyle`, `spotbugs`, `enforcer`, or `jacoco`. |
| `mavenPlugin` | string | Maven plugin associated with the finding. |
| `mavenPhase` | string | Maven phase most closely associated with the finding. |
| `testClass` | string | Test class, source file, or plugin log source depending on `sourceReportFormat`. |
| `testName` | string | Test method, source location, or plugin execution depending on `sourceReportFormat`. |
| `failureKind` | string | Normalized finding kind, such as `failure`, `error`, `violation`, `bug`, `rule`, or `threshold`. |
| `failureType` | string | Optional source-specific type, rule, category, or threshold metric. |
| `message` | string | Optional one-line finding message from the source artifact. |
| `reproduceCommand` | string | Minimal Maven command for reproducing or re-running the relevant check locally. |
| `confidence` | string | Confidence level for the finding. Current values use `high` for fixture-backed deterministic evidence. |
| `confidenceReasons` | array of string | Human-readable evidence explaining why the finding is trusted. |
| `sourceReportFormat` | string | Original source format, such as `junit-xml`, `checkstyle-xml`, `spotbugs-xml`, or `maven-log`. |

## Demo JSON

Generated from:

```bash
prmaven why -project demo/multi-module-failure -format json
```

```json
{
  "projectRoot": "<PROJECT_ROOT>/demo/multi-module-failure",
  "summary": {
    "moduleCount": 3,
    "reportCount": 2,
    "findingCount": 2
  },
  "modules": [
    {
      "name": "multi-module-failure",
      "path": ".",
      "pom": "pom.xml"
    },
    {
      "name": "payment-api",
      "path": "payment-api",
      "pom": "payment-api/pom.xml"
    },
    {
      "name": "payment-core",
      "path": "payment-core",
      "pom": "payment-core/pom.xml"
    }
  ],
  "findings": [
    {
      "id": "payment-api-paymentapiit-shouldexposepaymentsummary-error",
      "module": "payment-api",
      "modulePath": "payment-api",
      "reportPath": "payment-api/target/failsafe-reports/TEST-dev.prmaven.demo.PaymentApiIT.xml",
      "reportKind": "failsafe",
      "mavenPlugin": "maven-failsafe-plugin",
      "mavenPhase": "verify",
      "testClass": "dev.prmaven.demo.PaymentApiIT",
      "testName": "shouldExposePaymentSummary",
      "failureKind": "error",
      "failureType": "java.net.ConnectException",
      "message": "java.net.ConnectException: connection refused",
      "reproduceCommand": "mvn -pl payment-api -am -Dit.test=PaymentApiIT verify",
      "confidence": "high",
      "confidenceReasons": [
        "failure was found in a Maven Failsafe JUnit XML report",
        "report path maps to Maven module payment-api",
        "reproduction command targets test class PaymentApiIT"
      ],
      "sourceReportFormat": "junit-xml"
    },
    {
      "id": "payment-core-paymentroundingtest-shouldrejectinvalidscale-failure",
      "module": "payment-core",
      "modulePath": "payment-core",
      "reportPath": "payment-core/target/surefire-reports/TEST-dev.prmaven.demo.PaymentRoundingTest.xml",
      "reportKind": "surefire",
      "mavenPlugin": "maven-surefire-plugin",
      "mavenPhase": "test",
      "testClass": "dev.prmaven.demo.PaymentRoundingTest",
      "testName": "shouldRejectInvalidScale",
      "failureKind": "failure",
      "failureType": "org.opentest4j.AssertionFailedError",
      "message": "expected:\u003c2.50\u003e but was:\u003c2.49\u003e",
      "reproduceCommand": "mvn -pl payment-core -am -Dtest=PaymentRoundingTest test",
      "confidence": "high",
      "confidenceReasons": [
        "failure was found in a Maven Surefire JUnit XML report",
        "report path maps to Maven module payment-core",
        "reproduction command targets test class PaymentRoundingTest"
      ],
      "sourceReportFormat": "junit-xml"
    }
  ]
}
```
