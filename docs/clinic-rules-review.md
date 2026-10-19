# Clinic rules review — DEC01 and DEC02

This is the review record for `TASK-SET-02` (Member B; clinic reviewer). It turns the open decisions in SRS v1.0 into questions that can be answered and signed off. Member B has confirmed their role, but their personal name has not been supplied. The review packet is prepared on `member-b/confirm-clinic-rules`. **No clinic-specific answer or approval is recorded yet.** The requirements below are the SRS baseline; they do not substitute for a clinic decision. Keep `docs/decisions.md` marked pending until the relevant reviewers approve a dated answer.

## Review record

| Item | Value to record |
| --- | --- |
| Project and institution | AyurFlow — Group 04, IS5109 Community Project; Sabaragamuwa University of Sri Lanka, Faculty of Computing, Department of Computing and Information Systems (SRS v1.0 document control) |
| Clinic reviewer name and role | Pending |
| Team Member B name | Role confirmed by the current user; name pending |
| Team lead / supervisor names and roles | Team lead pending. SRS v1.0 names Mrs K. T. Dananjali as supervisor, with spelling and designation still to be confirmed. |
| Review date and SRS version | Pending; SRS v1.0 is the current review baseline |
| Evidence location (signed notes, letter, or meeting record) | Pending; store approved evidence in the team's agreed location |
| Outcome | Pending — approved, approved with conditions, or returned for revision |

Do not put real patient records, passwords, or private contact details in this repository. Record the approved rule and a reference to its review evidence here.

## Working rules for development

The current user has directed this planning task to proceed without waiting for clinic signoff. The rules in this section are a **provisional development contract**, not clinic answers. They allow synthetic-data schema and interface work to proceed. The pending questions and named gates below remain open for clinical UAT and any real-data pilot. If a clinic answer differs, update the contract, migrations, screens and affected tests before acceptance.

| Topic | Interim rule for synthetic development | Open gate and accountable owner |
| --- | --- | --- |
| Clinic identity | Use the academic project identity above. Display “Partner clinic” in internal demos; do not put an invented clinic name on records or release material. | Team lead obtains clinic legal/display name, organisation letter, authorised reviewer and dates under DEC01 before SRS/clinic review. |
| Patient registration | Require a server-generated patient ID and a nonblank Unicode name. Treat phone as optional in synthetic examples. Do not require NIC/passport, birth date, address or other demographics. | B and clinic clinician confirm the allowed/required contact fields, labels, formats and limits before registration UAT. |
| Supplied identity | Store an optional identity type and original value separately from a canonical value. For synthetic examples, trim outer whitespace, remove ASCII spaces and hyphens, and uppercase ASCII letters; uniqueness is on type plus canonical value when supplied. Do not merge NIC and passport namespaces. Empty canonical values are rejected. | B and clinic clinician review equivalent/distinct NIC and passport examples, other accepted types and punctuation/country rules before patient identity acceptance. |
| Clinical record access | Deny full notes to non-clinical roles. A clinician creating a draft may access it; another clinician needs an explicit patient-record assignment. Deny access when no assignment exists. A disabled account cannot act. Reception may use permitted demographics and bookings; an assigned therapist may use only the session details and observations needed for delivery; dispensary and billing may use the minimum patient/session identifiers needed for their work. | B and clinic clinician confirm who grants/changes assignments, historical access and therapist observation scope under DEC02; A reviews the backend enforcement contract. |
| Consultation history | Permit a draft without a plan. Signing freezes its original text and records signer/time. Corrections are separate dated addenda with author and reason. | B and clinic clinician confirm draft/signing fields and addendum reason rules before consultation UAT. |
| Single pathway | Use the SRS's provisional Preparatory → Primary → Recovery order only in synthetic examples. A clinician approves the current stage with actor/time; no stage skip or automatic approval. For synthetic transitions, a clinician records stage completion after all non-cancelled sessions for that stage are completed; a booked or in-progress session keeps it open. Advancement is a separate clinician action. After the final stage is completed the plan is read-only for advancement and new starts. Session start checks active plan, current approved stage and assigned therapist. Billing status has no effect. | B, C and clinic clinician confirm labels, approval and completion checklists, including whether a stage with no sessions may complete, before clinical UAT. |
| Hold and cancellation | Hold blocks new starts and advancement while retaining approvals. Cancel is terminal for the pilot, cancels future unstarted sessions and preserves completed history. An already started session can record completion during hold. | B, C and clinic reviewer confirm DEC08 races and in-progress behavior before SRS approval. |

The interim identity rule is deliberately narrow: it makes synthetic duplicate examples testable without asserting a clinic-approved national-identifier format. The API and database must preserve the original supplied value and avoid disclosing another patient when a duplicate is detected.

Synthetic normalization examples: `NIC: 1999 12345678 V` and `NIC: 1999-12345678v` both become `NIC:199912345678V`; `Passport: ab-12345` becomes `Passport:AB12345`. The same canonical characters under different identity types remain distinct. A patient with no supplied identity has no canonical value and may still be registered. These examples are for development only; the clinic must confirm the identity types and equivalences before AT02 acceptance.

## DEC01 — clinic and academic identity (week 1 gate)

| Question to confirm with clinic, team, and supervisor | Recorded answer / evidence |
| --- | --- |
| What is the clinic's approved name and the organisation named in the proposal and pilot materials? | Pending |
| Who may speak for the clinic in clinical workflow review and pilot acceptance? Record role and approval authority. | Pending |
| Is the organisation letter available, and where is the approved copy kept? | Pending |
| Confirm the academic supervisor's spelling, designation, and required review or approval. | SRS v1.0 provisionally names Mrs K. T. Dananjali; confirmation pending. |
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

| Decision | Current state | Needed to close |
| --- | --- | --- |
| DEC01 | Pending clinic/team/supervisor confirmation | Clinic identity, authorised reviewer, organisation letter, supervisor details, review dates and evidence reference |
| DEC02 | Pending clinic clinician confirmation | Patient/consultation fields, identity normalisation, record access, pathway order, approval and completion rules, with examples and evidence reference |

Before marking `TASK-SET-02` complete, record dated answers and clinic evidence for DEC01–02, obtain the required clinic/supervisor review, and update `docs/decisions.md`. Then give A the agreed fields and access rules for the shared schema/auth contract, C the stage/session clearance examples, and F the screen states. Record each handoff's recipient, date and reviewed contract link here. B can draft migrations and UI against reviewed examples, but affected feature acceptance remains gated on approved answers.

| Handoff | Recipient/date | Reviewed evidence |
| --- | --- | --- |
| Agreed patient and consultation fields, identity handling and clinical record access | A — pending | Pending |
| Agreed stages, approval/completion and session-start examples, including DEC08 interaction | C — pending | Pending |
| Registration, consultation and plan screen states and allowed actions | F — pending | Pending |

Sources: AyurFlow Software Requirements Specification v1.0, sections 2, 5–7, 12–18; AyurFlow Team Task Plan, `TASK-SET-02`, `TASK-SET-06`, `TASK-UI-02`, and Member B ownership.
