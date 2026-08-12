# Hardware Device Ideas

This is a catalog of possible constrained-device applications for PromiseGrid.
It is intentionally exploratory: none of these entries selects a pCID, freezes
a payload, chooses hardware, or commits the repository to a product. The point
is to choose a real physical job first, then design the smallest honest device
promise around it.

## What makes a device idea Grid-shaped

A useful hardware participant should make a narrow, durable statement about
what it can actually observe or do. It should sign and retain or queue its own
exact records within its limits; use a carrier only to move bytes; and let
workers, supervisors, and other local applications make their own later
assessments. A barcode scan, button press, sensor reading, Wi-Fi connection,
or dashboard login alone is not a promise made by a person.

For every idea below, a future design would need to state the device key and
replacement story, flash/power-loss behavior, offline queue rules, duplicate
handling, clock limits, physical-tampering assumptions, chosen frozen specs,
and every responsibility delegated to a companion agent or carrier. This is
the constrained-runtime discipline described by the PromiseGrid development
guide; it does not require the device to host every possible layer.

## Ideas

### 1. Handheld equipment inventory scanner

**Physical job.** A worker scans a tagged tool, fixture, meter, laptop, or
other durable asset while walking the floor, shelves, vehicles, or cages.

**Narrow device record.** “This device scanned asset identifier X in inventory
context Y at local sequence N.” The context may be a selected room, shelf,
cart, or inventory round.

**Why it belongs on the Grid.** The scan is durable evidence that can survive
an offline walk-through and later be assessed by a supervisor without relying
on one hosted inventory database to be the only memory of the count.

**Limits.** It does not prove ownership, custody, condition, quantity beyond
the scan rule, or that the item stayed there after the observation.

### 2. Cycle-count station

**Physical job.** A fixed or mobile station counts stock in a bin, rack, or
pallet location using barcode scans and a small keypad or buttons.

**Narrow device record.** “This device completed count round R for location L
and observed item/SKU S with entered count Q.”

**Why it belongs on the Grid.** A later reconciliation can compare several
independent count observations, their device histories, and their exact input
records instead of overwriting one mutable number.

**Limits.** It does not establish the true quantity, authorize adjustment, or
prove that a worker counted correctly.

### 3. Receiving-dock label verifier

**Physical job.** At receiving, a worker scans a purchase-order, shipment,
carton, or lot label as goods arrive.

**Narrow device record.** “This device observed label set A at dock/location D
during receiving session R.”

**Why it belongs on the Grid.** Receiving is a handoff across vendors,
carriers, workers, and supervisors. Exact signed observations preserve what
the device saw even if another party later disputes a shipment.

**Limits.** It does not prove that the contents match the label, that goods are
undamaged, or that a supplier fulfilled a contract.

### 4. Shelf and location audit wand

**Physical job.** A small scanner verifies that an item label was seen at a
specific shelf, aisle, vehicle, or staging area.

**Narrow device record.** “This device observed asset/item X together with
location marker L during audit round R.”

**Why it belongs on the Grid.** It gives a portable, durable trail of local
location observations while allowing each operation to choose its own
reconciliation policy.

**Limits.** It does not provide continuous tracking or prove that the object
was not moved immediately after scanning.

### 5. Tool-return condition kiosk

**Physical job.** A worker scans a returned shared tool and selects a simple
condition such as ready, damaged, missing-part, or needs-review.

**Narrow device record.** “This device recorded a return inspection input for
tool X with stated condition C.”

**Why it belongs on the Grid.** The physical return point emits a durable,
reviewable observation that maintenance or supervision can assess later rather
than overwriting the previous state silently.

**Limits.** The device does not diagnose the tool, assign blame, or decide
whether the condition is accurate.

### 6. Maintenance-meter capture device

**Physical job.** A technician scans an asset tag and records a meter reading
from a machine, compressor, generator, lift, or production counter.

**Narrow device record.** “This device observed displayed meter value M for
asset X during maintenance round R.”

**Why it belongs on the Grid.** Exact observations can be retained through
offline maintenance rounds and assessed later against other maintenance
evidence without turning the device into a maintenance authority.

**Limits.** It does not prove the meter is calibrated, the value is truthful,
or that any maintenance action occurred.

### 7. Temperature and humidity logger

**Physical job.** A low-power device periodically measures an environmental
condition in a storage room, cabinet, vehicle, or work area.

**Narrow device record.** “This device measured temperature T and humidity H
at its local sample sequence N.”

**Why it belongs on the Grid.** It is a constrained participant with useful
offline evidence: it can queue exact observations and later carry them to
interested parties without claiming that a dashboard is the record of truth.

**Limits.** It does not prove conditions everywhere in the room, product
safety, calibration, or legal compliance.

### 8. Cold-room door monitor

**Physical job.** A battery device senses door-open and door-closed transitions
on a cooler, freezer, secured cabinet, or warehouse bay.

**Narrow device record.** “This device observed its door sensor transition to
state S at local sequence N.”

**Why it belongs on the Grid.** Parties can independently interpret a durable
transition history for operations or maintenance, even when connectivity is
intermittent.

**Limits.** It does not know who opened the door, why, whether it sealed,
or whether stored goods were affected.

### 9. Package handoff scanner

**Physical job.** A worker scans a package, tote, or pallet at a handoff point
between work areas, shifts, carriers, or outbound staging.

**Narrow device record.** “This device observed package identifier P at handoff
marker H during round R.”

**Why it belongs on the Grid.** The record is an independently transportable
observation that can be compared by workers and supervisors without asserting
a universal custody service.

**Limits.** It does not prove possession, delivery, contents, or acceptance by
another person.

### 10. Connected scale station

**Physical job.** A scale measures a box, ingredient, component batch, or scrap
container, with an optional scanned item identifier.

**Narrow device record.** “This device measured weight W while item context X
was selected.”

**Why it belongs on the Grid.** A signed measurement record preserves the
device’s local evidence for later quality or inventory policy, including when
the scale is temporarily disconnected.

**Limits.** It does not prove the item on the scale, calibration, tare policy,
or that the measurement authorizes any transaction.

### 11. Damage and incident marker

**Physical job.** A rugged button-and-scan device lets a worker mark a visible
damage, spill, blocked aisle, or equipment incident at the point of discovery.

**Narrow device record.** “This device recorded incident category C for scanned
context X, with optional locally captured measurement or image reference.”

**Why it belongs on the Grid.** It retains a durable local report that can be
carried and reviewed by different applications without forcing an immediate
global conclusion.

**Limits.** It does not establish cause, severity, liability, or remedy.

### 12. Safety-inspection checkpoint

**Physical job.** A technician works through a short, repeatable inspection at
a ladder rack, lift, exit, machine guard, eyewash station, or vehicle.

**Narrow device record.** “This device recorded the selected inspection inputs
for checklist instance I at asset/location X.”

**Why it belongs on the Grid.** The inspection inputs remain exact evidence;
each organization can independently decide what they mean and what follow-up
they require.

**Limits.** It does not certify safety, enforce a rule, or replace professional
inspection judgment.

### 13. Calibration-bench reader

**Physical job.** A bench connects to a meter, gauge, sensor, or tester and
captures a measurement against a scanned instrument identifier.

**Narrow device record.** “This device received measurement M from interface I
while instrument X was selected.”

**Why it belongs on the Grid.** It separates durable source measurement from
the later local decision that a device passed, failed, or may remain in use.

**Limits.** It does not establish traceability, calibration validity, or an
approval to release the instrument.

### 14. Work-cell start/stop marker

**Physical job.** A button panel or scanner marks a local work cell’s setup,
start, pause, completion, or cleanup activity.

**Narrow device record.** “This device observed operator input event E at work
cell C during local sequence N.”

**Why it belongs on the Grid.** A compact physical participant can retain
production-floor observations through outages without redefining them as an
authoritative production schedule.

**Limits.** It does not prove work was performed, quantity produced, quality,
or employee attendance.

### 15. Machine changeover verifier

**Physical job.** During a line or machine changeover, a worker scans the old
and new material/tooling labels and confirms selected visible steps.

**Narrow device record.** “This device observed selected labels and entered
changeover inputs for machine M and round R.”

**Why it belongs on the Grid.** It provides reviewable evidence across shifts
and roles while keeping approval, release, and quality policy outside the
scanner’s narrow claim.

**Limits.** It does not prove a changeover was correctly performed or authorize
production.

### 16. Scrap and rework bin scanner

**Physical job.** A worker scans parts placed in a scrap or rework container
and selects a local reason category.

**Narrow device record.** “This device observed item/lot X placed at bin B with
entered category C.”

**Why it belongs on the Grid.** It preserves the original floor observation so
later quality, inventory, and supervision views need not collapse it into one
mutable status.

**Limits.** It does not determine root cause, financial disposition, or whether
the part actually stayed in the bin.

### 17. Rental or loan return station

**Physical job.** A compact station scans an item returned from a customer,
jobsite, department, or employee and records a quick visual input.

**Narrow device record.** “This device observed return scan X with entered
return condition C at station S.”

**Why it belongs on the Grid.** A return is a disputed physical moment. The
station’s exact record can travel with later review without asserting that it
settles ownership, payment, or damage responsibility.

**Limits.** It does not prove who returned it, contract terms, or chargeability.

### 18. Field-service kit auditor

**Physical job.** A technician scans the contents of a mobile service kit,
vehicle tote, emergency cabinet, or installation case before or after a job.

**Narrow device record.** “This device observed the listed identifiers during
kit audit K at local sequence N.”

**Why it belongs on the Grid.** It makes an offline, portable audit trail that
can later be reconciled by the company without requiring live warehouse access
at the jobsite.

**Limits.** It does not prove completeness, consumption, or possession between
audits.

### 19. Lot and expiry label reader

**Physical job.** A scanner reads lot, serial, and expiry labels during
receiving, picking, replenishment, or quality checks.

**Narrow device record.** “This device observed encoded lot/expiry label L in
operation context O.”

**Why it belongs on the Grid.** Exact original label evidence can be retained
and assessed by different local policies, including when an operation cannot
reach its usual software.

**Limits.** It does not prove the label is correct, the item matches it, or that
the item is permitted for use.

### 20. Backup-power and equipment-runtime logger

**Physical job.** A small device records power transitions, battery state, or
runtime counters for a pump, freezer, network cabinet, generator, or machine.

**Narrow device record.** “This device observed power/runtime measurement M at
local sequence N.”

**Why it belongs on the Grid.** It is a natural constrained participant: it can
retain exact, queueable operational observations through the same outages it is
measuring, and later carry them without turning a carrier into authority.

**Limits.** It does not establish downtime cause, safety, service compliance,
or that the monitored equipment behaved correctly.

## Choosing one later

The best next example is likely the one with a single physical action, a
credible offline period, a clearly bounded observation, and a real reason for
workers and supervisors to inspect the same exact evidence later. A handheld
inventory scanner is intentionally the smallest starting point in this list;
it can remain simply “take inventory” until a chosen scenario requires more.
