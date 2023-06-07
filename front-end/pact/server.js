import pact from "@pact-foundation/pact-core";
export const server = pact.createServer({ port: 9999 });
server.start().then(() => {
    console.log('Server started');
    server.delete();
    server.stop();
});