# Patients module

Proposed owner: B. Reviewer: C. SRS: PAT-01-02; CON-01-03.

Patient registration/search with optional identity, clinician consultation drafts, immutable signing and dated addenda.

Status: TASK-PAT-01 schema and provisional validation are drafted on `member-b/patient-schema`. Migration `000005_patients.up.sql` adds patient demographics and optional supplied identity with database-generated canonical value and a unique type/canonical constraint. A must confirm the migration number; C reviews the schema; DEC02 still controls final fields and identity equivalence. `validation.go` mirrors the provisional normalization rule. No registration/search route is registered, so AT02 remains pending.

PAT-02 will add handler/service/repository behavior with parameterized pgx queries, backend role and record checks, CSRF, transactional audit, bounded search, and duplicate-race handling. Do not register protected routes before those controls exist.
