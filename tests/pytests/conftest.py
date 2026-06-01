# Copyright © 2025-2026 OpenCHAMI a Series of LF Projects, LLC
#
# SPDX-License-Identifier: MIT

import pytest
import requests
import json
import time
import pprint
import os
import glob


def pytest_addoption(parser):
    parser.addoption(
        "--enable-redfish-discovery",
        action="store_true",
        default=False,
        help="Enable hardware discovery during test setup",
    )
    parser.addoption(
        "--discovery-json",
        action="store",
        default=None,
        metavar="DIR",
        help="Directory containing JSON files to POST to /v2/Inventory/RedfishEndpoints",
    )

smd_base_url = "http://smd:27779/hsm"
inventory_base_url = "http://inventory:8080/hsm"


@pytest.fixture(scope="session")
def discover_hardware(request):
    # setup
    is_redfish_discovery_enabled = request.config.getoption("--enable-redfish-discovery")
    discovery_json_dir = request.config.getoption("--discovery-json")

    if discovery_json_dir is not None:
        json_files = glob.glob(os.path.join(discovery_json_dir, "*.json"))
        if not json_files:
            pytest.fail(f"No JSON files found in {discovery_json_dir}")
        for json_file in json_files:
            print(f"POSTing {json_file} to RedfishEndpoints")
            with open(json_file) as f:
                payload = json.load(f)
            response = requests.post(f"{smd_base_url}/v2/Inventory/RedfishEndpoints", json=payload)
            if not response.ok:
                response_failure("POST", response)

    if is_redfish_discovery_enabled:
        bmc_nodes = [ "x0c0s1b0", "x0c0s2b0", "x0c0s3b0", "x0c0s4b0" ]
        response = requests.get(f"{smd_base_url}/v2/Inventory/RedfishEndpoints")
        if not response.ok:
            print_response("GET", response)
            pytest.fail(f"Failed to get {response.url}")

        discovered_nodes = [ endpoint.get("ID")
                            for endpoint in json.loads(response.text).get("RedfishEndpoints", [])
                            if endpoint.get("DiscoveryInfo", {}).get("LastDiscoveryStatus") == "DiscoverOK"]
        undiscovered_nodes = list(set(bmc_nodes) - set(discovered_nodes))
        print(f"bmc_nodes: {bmc_nodes}")
        print(f"discovered_nodes: {discovered_nodes}")
        print(f"undiscovered_nodes: {undiscovered_nodes}")


        for node in undiscovered_nodes:
            print(f"discover: {node}")
            request_body = {
                    "RedfishEndpoints" : [
                        {
                         "ID" : node,
                         "FQDN" : node,
                         "RediscoverOnUpdate" : True,
                         "User" : "root",
                         "Password" : "root_password"
                         }]
                    }
            response = requests.post(f"{smd_base_url}/v2/Inventory/RedfishEndpoints", json=request_body)
            if not response.ok:
                print_response("POST", response)

        if undiscovered_nodes:
            for i in range(0, 10):
                done = True
                print(f"Waiting for discovery to finish. {i}")
                response = requests.get(f"{smd_base_url}/v2/Inventory/RedfishEndpoints")
                if response.ok:
                    endpoints = json.loads(response.text)
                    discovery_results = { endpoint.get("ID"): endpoint.get("DiscoveryInfo").get("LastDiscoveryStatus")
                                          for endpoint in endpoints.get("RedfishEndpoints")}
                    pprint.pprint(discovery_results)
                    for node in undiscovered_nodes:
                        endpoint = discovery_results.get(node)
                        print(f"{node} {endpoint}")
                        if endpoint != "DiscoverOK":
                            print(f"- {node} {endpoint}")
                            done = False
                    if done:
                        break
                time.sleep(1)

    replicate_components()
    replicate_component_endpoints()
    replicate_ethernet_interfaces()
    replicate_redfish_endpoints()
    replicate_service_endpoints()
    replicate_hardware()

    yield

    # tear down


def replicate_components():
    response = requests.get(f"{smd_base_url}/v2/State/Components")
    if not response.ok:
        print_response("GET", response)
    smd_components = json.loads(response.text)

    print("POST Components to the inventory service")
    response = requests.post(f"{inventory_base_url}/v2/State/Components", json=smd_components)
    if not response.ok:
        print_response("POST", response)


def replicate_component_endpoints():
    response = requests.get(f"{smd_base_url}/v2/Inventory/ComponentEndpoints")
    if not response.ok:
        print_response("GET", response)
    smd_components = json.loads(response.text)

    print("POST ComponentEndpoints to the inventory service")
    response = requests.post(f"{inventory_base_url}/v2/Inventory/ComponentEndpoints", json=smd_components)
    if not response.ok:
        print_response("POST", response)


def replicate_ethernet_interfaces():
    response = requests.get(f"{smd_base_url}/v2/Inventory/EthernetInterfaces")
    if not response.ok:
        print_response("GET", response)
    ethernet_interfaces = json.loads(response.text)

    print("POST EthernetInterfaces to the inventory service")
    for eth in ethernet_interfaces:
        response = requests.post(f"{inventory_base_url}/v2/Inventory/EthernetInterfaces", json=eth)
        if not response.ok:
            print_response("POST", response)


def replicate_redfish_endpoints():
    response = requests.get(f"{smd_base_url}/v2/Inventory/RedfishEndpoints")
    if not response.ok:
        print_response("GET", response)
    redfish_endpoints  = json.loads(response.text)

    print("POST RedfishEndpoints to the inventory service")
    for redfish_endpoint in redfish_endpoints.get("RedfishEndpoints"):
        response = requests.post(f"{inventory_base_url}/v2/Inventory/RedfishEndpoints", json=redfish_endpoint)
        if not response.ok:
            print_response("POST", response)


def replicate_service_endpoints():
    response = requests.get(f"{smd_base_url}/v2/Inventory/ServiceEndpoints")
    if not response.ok:
        print_response("GET", response)
    smd_service_endpoints = json.loads(response.text)

    print("POST ServiceEndpoints to the inventory service")
    response = requests.post(f"{inventory_base_url}/v2/Inventory/ServiceEndpoints", json=smd_service_endpoints)
    if not response.ok:
        print_response("POST", response)


def replicate_hardware():
    response = requests.get(f"{smd_base_url}/v2/Inventory/Hardware")
    if not response.ok:
        print_response("GET", response)
    smd_hardware = json.loads(response.text)

    print("POST Hardware to the inventory service")
    smd_hardware_post_obj = { "Hardware" : smd_hardware }
    response = requests.post(f"{inventory_base_url}/v2/Inventory/Hardware", json=smd_hardware_post_obj)
    if not response.ok:
        print_response("POST", response)


def print_response(method, response):
        print(f"{method} URL: {response.url}, Code: {response.status_code}, Body:")
        print(response.text)
        print(json.dumps(response.text, indent=4))
