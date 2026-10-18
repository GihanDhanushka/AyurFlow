# Clinic rules review — DEC01 and DEC02

This is the review record for `TASK-SET-02` (Member B; clinic reviewer). It turns the open decisions in SRS v1.0 into questions that can be answered and signed off. **No clinic-specific answer or approval is recorded yet.** The requirements below are the SRS baseline; they do not substitute for a clinic decision. Keep `docs/decisions.md` marked pending until the relevant reviewers approve a dated answer.

## Review record

| Item | Value to record |
| --- | --- |
| Clinic reviewer name and role | Pending |
| Team Member B name | Pending |
| Team lead / supervisor names and roles | Pending |
| Review date and SRS version | Pending; SRS v1.0 is the current review baseline |
| Evidence location (signed notes, letter, or meeting record) | Pending; store approved evidence in the team's agreed location |
| Outcome | Pending — approved, approved with conditions, or returned for revision |

Do not put real patient records, passwords, or private contact details in this repository. Record the approved rule and a reference to its review evidence here.

## DEC01 — clinic and academic identity (week 1 gate)

| Question to confirm with clinic, team, and supervisor | Recorded answer / evidence |
| --- | --- |
| What is the clinic's approved name and the organisation named in the proposal and pilot materials? | Pending |
| Who may speak for the clinic in clinical workflow review and pilot acceptance? Record role and approval authority. | Pending |
| Is the organisation letter available, and where is the approved copy kept? | Pending |
| Who is the academic supervisor, and what review or approval is required from them? | Pending |
| What are the kickoff, requirements review, clinic demonstration, and acceptance review dates? | Pending |

DEC01 is complete only when the team, clinic, and supervisor agree these details and the review evidence is recorded. This does not itself authorise real patient data use; that is DEC05.

## DEC02 — patient and consultation fields (weeks 1–2 gate)

| Question for the clinic clinician and Member B | Recorded rule / example |
| --- | --- |
| Which contact and demographic fields are required at registration, optional, or prohibited? Specify labels, formats, maximum lengths, and whether blank values are allowed. | Pending |
| What identity types are accepted (NIC, passport, other), and what exact normalisation applies before duplicate checking (case, spaces, punctuation, country/type)? Provide equivalent and distinct examples. | Pending |
| Can one patient have more than one supplied identity, and how should an apparent duplicate be resolved by staff without exposing another patient's record? | Pending |
| Which fields may a receptionist change after registration, and which changes require a reason or clinician action? | Pending |
| Which consultation fields are required for a draft, and which must be present before signing? Include any structured observations versus free text. | Pending |
| What reason categories or minimum text are required for a signed-note addendum? | Pending |

The SRS already requires a server-generated patient ID, a name, optional NIC/passport, safe uniqueness for a *supplied normalised* identity, Unicode names, patient search, consultation drafts without a plan, immutable signed text, and separate dated addenda. The open answers above refine those rules; they cannot make NIC mandatory or allow rewriting a signed original without an approved scope change.

## DEC02 — clinical record access

| Question for the clinic clinician and Member B | Recorded rule / example |
| --- | --- |
| Which clinicians may read and edit a given patient's consultations and care plans: all assigned clinicians, a named care team, or another explicit policy? How is assignment created and changed? | Pending |
| What minimum patient and booking details may receptionists see? Confirm that consultation text stays hidden. | Pending |
| What session observations may the assigned therapist read or write, and can they see observations from other sessions? | Pending |
| Which patient/session identifiers may dispensary and billing staff see for their work? | Pending |
| Who may view historical clinical decisions and addenda after a clinician leaves or an account is disabled? | Pending |
| How are access-denied and inaccessible-record cases presented to staff without revealing another patient's details? | Pending |

Backend checks must enforce both role and patient-record access on direct API requests. An administrator-only account does not acquire clinical-note access. The clinic's answers must be reflected in the API contract and AT01/AT03/AT13 tests.

## DEC02 — one pathway, stages, and completion

| Question for the clinic clinician, Member B, and Member C | Recorded rule / example |
| --- | --- |
| What is the exact clinic-approved ordered pathway? Record each stage's label and order; confirm whether the SRS's provisional Preparatory → Primary → Recovery sequence is accepted or replaced. | Pending |
| What fields or checklist must a clinician complete to approve each stage? Is a reason required for approval or rejection? | Pending |
| What conditions make a stage complete, including the role of delivered sessions and observations? | Pending |
| What conditions make the final stage and plan complete and read-only for further advancement and new starts? | Pending |
| What plan/stage/session examples demonstrate that a booked session may start, and what examples must be rejected? | Pending |
| How should the clinic handle a session already in progress when its plan is put on hold? Confirm with DEC08 and Member C. | Pending |

The SRS requires ordered advancement, a named clinician's recorded approval, no stage skipping, and clearance **at session start** using the current stage, plan status, and required approval. Booking does not grant clearance, and invoice balance does not affect it. A hold blocks new starts and advancement without erasing prior approvals. Cancellation preserves history and cancels future unstarted sessions. DEC08 separately confirms the finer hold, cancellation, and concurrency semantics.

## Examples to review before approval

Use synthetic records only. Add clinic-specific expected outcomes to the evidence record.

| Example | SRS baseline expected result | Clinic detail still needed |
| --- | --- | --- |
| Register a patient without NIC or passport. | Registration succeeds with an internal ID. | Required contact fields. |
| Two staff submit the same supplied identity at once. | At most one patient is created for the same normalised identity; the other request reports a safe conflict. | Normalisation examples and duplicate-handling procedure. |
| Search for a patient with a Sinhala or other Unicode name. | Name round-trips and authorised search finds the record. | Agreed search behavior and record access. |
| Sign a consultation while another clinician edits the draft. | Signed original cannot be overwritten; a later correction is a dated addendum. | Required signing fields and correction reason. |
| Approve or start a session for a future stage. | Stage skip or noncurrent-stage start is rejected. | Exact pathway and clearance checklist. |
| Hold or cancel a plan while a session start request races. | Requests observe a consistent decision; no new start bypasses an effective hold or cancellation. | DEC08 review of in-progress delivery and decision order. |

## Closure and handoff

Before marking `TASK-SET-02` complete, record dated answers and clinic evidence for DEC01–02, obtain the required clinic/supervisor review, and update `docs/decisions.md`. Then give A the agreed fields and access rules for the shared schema/auth contract, C the stage/session clearance examples, and F the screen states. B can draft migrations and UI against reviewed examples, but affected feature acceptance remains gated on approved answers.

Sources: AyurFlow Software Requirements Specification v1.0, sections 2, 5–7, 12–18; AyurFlow Team Task Plan, `TASK-SET-02`, `TASK-SET-06`, `TASK-UI-02`, and Member B ownership.
