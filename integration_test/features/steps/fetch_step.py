import os
import sys
import time

from hamcrest import assert_that, is_, is_not
from behave.matchers import register_type
from utils import parse_optional


parse_optional.pattern = r'\s?\w*\s?'
register_type(optional=parse_optional)


@when('I fetch{external:optional}"{profile}" "{document}"')
def step_impl(context, external, profile, document):
    host = ""
    if external == "external":
        host = "workload"

    if document == "SVID":
        bin_file = "c_client"
    else:
        bin_file = "c_client_bundle"
    c_client_bin = os.popen("/mnt/c-spiffe/integration_test/helpers/bash-spire-scripts/ssh-fetch.sh %s %s %s %s" % (bin_file, document.lower(), profile.lower(), host))
    context.result = c_client_bin.read()


@then('I check that the "{document}" is returned correctly')
def step_impl(context, document):
    assert_that(context.result.find("error"), is_(-1), "There was an error")
    assert_that(context.result.find("Address: "), is_not(-1), "There is no Address")
    result = context.result.splitlines()
    document_content = result[0].split(" ")[-1]
    if document.lower() == "svid":
        context.svid = document_content
    else:
        context.bundle = document_content
    assert_that(document_content, is_not("(nil)"))


@then('I check that the "{document}" is not returned')
def step_impl(context, document):
    assert_that(context.result.find("error"), is_not(-1), "There was no error")
    assert_that(context.result.find("Address: "), is_not(-1), "There is no Address")
    result = context.result.splitlines()
    document_content = result[1].split(" ")[-1]
    if document.lower() == "svid":
        context.svid = document_content
    else:
        context.bundle = document_content
    assert_that(document_content, is_("(nil)"))


@when('The agent is turned off')
def step_impl(context):
    os.system("pkill spire-agent")
    time.sleep(5)


@when('The agent is turned on')
def step_impl(context):
    os.system("/mnt/c-spiffe/integration_test/helpers/bash-spire-scripts/grpc_connect_agent.sh")
    time.sleep(5)