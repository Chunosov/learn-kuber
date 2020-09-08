# Copyright 2015 gRPC authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
"""The Python implementation of the GRPC helloworld.Greeter client."""

from __future__ import print_function
import logging
import sys
import grpc
import helloworld_pb2
import helloworld_pb2_grpc


def run():
    # NOTE(gRPC Python Team): .close() is possible on a channel and should be
    # used in circumstances in which the with statement does not fit the needs
    # of the code.
    with grpc.insecure_channel(sys.argv[1],
        # NOTE: this doesn't work as well: https://github.com/improbable-eng/grpc-web/issues/266
        # NOTE: There is the list of available working:
        # NOTE: https://github.com/grpc/grpc/blob/v1.31.x/include/grpc/impl/codegen/grpc_types.h
        options=(('grpc.default_authority', 'grpc-greeter.default.example.com'),)) as channel:
        stub = helloworld_pb2_grpc.GreeterStub(channel)
        response, call = stub.SayHello.with_call(
            helloworld_pb2.HelloRequest(name='you'),
            # NOTE: We heed the "Host" header, but it doesn't support upper-case:
            # NOTE: https://github.com/grpc/grpc/issues/9863
            # NOTE: This is example with metadata:
            # NOTE: https://github.com/grpc/grpc/blob/master/examples/python/metadata/metadata_client.py
            #metadata=(('host', 'grpc-greeter.default.example.com'),)
            #metadata=(('host', 'grpc-greeter.default.svc.cluster.local'),)
            #metadata=(('authority', 'grpc-greeter.default.example.com'),)
            #metadata=(('authority', 'grpc-greeter.default.svc.cluster.local'),)
            # NOTE: There were two PRs in python grpc repo providing "authority" and "host" params
            # NOTE: but both of them were closed without merging.
            # NOTE: https://github.com/grpc/grpc/pull/14077
            # NOTE: https://github.com/grpc/grpc/pull/14361
            #authority='grpc-greeter.default.svc.cluster.local'
            #host='grpc-greeter.default.svc.cluster.local'
        )
        print("Greeter client received: " + response.message)


if __name__ == '__main__':
    logging.basicConfig()
    run()
