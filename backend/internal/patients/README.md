# Patients module

Proposed owner: B. Reviewer: C. SRS: PAT-01-02; CON-01-03.

Patient registration/search with optional identity, clinician consultation drafts, immutable signing and dated addenda.

Status: TASK-PAT-01 schema and provisional validation were merged through PR #19. Migration `000005_patients.up.sql` adds patient demographics and optional supplied identity with database-generated canonical value and a unique type/canonical constraint. A must confirm the migration number; C reviews the schema; DEC02 still controls final fields and identity equivalence. `validation.go` mirrors the provisional normalization rule.

TASK-PAT-02 work on `member-b/registration-search` adds a service with parameterized pgx queries for registration, ID/name/identity/contact search, demographic updates with version checks, and database-arbitrated duplicate identity conflicts. Patient creation and correction write structured audit events in the mutation transaction; a failed audit rolls back the change. Search returns at most 25 records after each candidate passes the supplied record-access policy. A missing policy denies every operation. The integration test exercises synthetic optional identities, a concurrent duplicate race, Unicode and other searches, access filtering, stale updates and forced audit failure.

The service deliberately has no HTTP registration. A's AUTH-03 role/record policy and mutation CSRF integration, A/C review of the shared contract, and DEC02 field/access decisions are still needed before protected routes and OpenAPI publication. PAT-03 owns the patient screens. AT02 remains pending until the route and screen are integrated and acceptance evidence is recorded.
