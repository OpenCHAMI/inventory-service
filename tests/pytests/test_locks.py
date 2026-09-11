# Copyright © 2026 OpenCHAMI a Series of LF Projects, LLC
#
# SPDX-License-Identifier: MIT

import pytest
import requests

from conftest import inventory_base_url, print_response

locks_base = f"{inventory_base_url}/v2/locks"


def _first_component_id():
    """Return the xname of an existing component, or None."""
    response = requests.get(f"{inventory_base_url}/v2/State/Components")
    if not response.ok:
        return None
    components = response.json().get("Components", [])
    if not components:
        return None
    return components[0].get("ID")


def _status_for(xname):
    """Return the /locks/status entry for a single xname, or None."""
    response = requests.post(f"{locks_base}/status", json={"ComponentIDs": [xname]})
    if not response.ok:
        print_response("POST", response)
        pytest.fail(f"Failed to POST /locks/status: {response.status_code}")
    for comp in response.json().get("Components", []):
        if comp.get("ID") == xname:
            return comp
    return None


def test_lock_unlock(discover_hardware):
    xname = _first_component_id()
    if xname is None:
        pytest.skip("no components available to lock")

    # Lock the component.
    response = requests.post(f"{locks_base}/lock", json={"ComponentIDs": [xname]})
    if not response.ok:
        print_response("POST", response)
        pytest.fail(f"Failed to POST /locks/lock: {response.status_code}")
    counts = response.json().get("Counts", {})
    assert counts.get("Success") == 1, f"expected 1 lock success, got {counts}"

    status = _status_for(xname)
    assert status is not None and status.get("Locked") is True, \
        f"expected component {xname} to be Locked, got {status}"

    # Unlock the component.
    response = requests.post(f"{locks_base}/unlock", json={"ComponentIDs": [xname]})
    if not response.ok:
        print_response("POST", response)
        pytest.fail(f"Failed to POST /locks/unlock: {response.status_code}")

    status = _status_for(xname)
    assert status is not None and status.get("Locked") is False, \
        f"expected component {xname} to be Unlocked, got {status}"


def test_reservation_lifecycle(discover_hardware):
    xname = _first_component_id()
    if xname is None:
        pytest.skip("no components available to reserve")

    # Make sure the component is unlocked and reservations are enabled.
    requests.post(f"{locks_base}/unlock", json={"ComponentIDs": [xname]})
    requests.post(f"{locks_base}/repair", json={"ComponentIDs": [xname]})

    # Create an admin reservation.
    response = requests.post(f"{locks_base}/reservations", json={"ComponentIDs": [xname]})
    if not response.ok:
        print_response("POST", response)
        pytest.fail(f"Failed to POST /locks/reservations: {response.status_code}")
    body = response.json()
    successes = body.get("Success", [])
    assert len(successes) == 1, f"expected 1 reservation, got {body}"
    reservation = successes[0]
    reservation_key = reservation.get("ReservationKey")
    assert reservation_key, f"expected a reservation key, got {reservation}"

    status = _status_for(xname)
    assert status is not None and status.get("Reserved") is True, \
        f"expected component {xname} to be Reserved, got {status}"

    # Release the reservation using its key.
    response = requests.post(
        f"{locks_base}/reservations/release",
        json={"ReservationKeys": [{"ID": xname, "Key": reservation_key}]},
    )
    if not response.ok:
        print_response("POST", response)
        pytest.fail(f"Failed to POST /locks/reservations/release: {response.status_code}")
    counts = response.json().get("Counts", {})
    assert counts.get("Success") == 1, f"expected 1 release success, got {counts}"

    status = _status_for(xname)
    assert status is not None and status.get("Reserved") is False, \
        f"expected component {xname} to be un-reserved, got {status}"


def test_status_not_found(discover_hardware):
    response = requests.post(f"{locks_base}/status", json={"ComponentIDs": ["xDOESNOTEXIST0"]})
    if not response.ok:
        print_response("POST", response)
        pytest.fail(f"Failed to POST /locks/status: {response.status_code}")
    assert "xDOESNOTEXIST0" in response.json().get("NotFound", []), \
        f"expected unknown xname in NotFound, got {response.text}"
