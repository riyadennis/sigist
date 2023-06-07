import { Pact } from '@pact-foundation/pact'
import path from 'path'


export const provider = new Pact({
   port: 20002,
   log: path.resolve(process.cwd(), 'logs', 'mockserver-integration.log'),
   dir: path.resolve(process.cwd(), 'pacts'),
   pactfileWriteMode: 'overwrite',
   consumer: 'GraphQLConsumer',
   provider: 'GraphQLProvider'
})

beforeAll(() => provider.setup())
afterAll(() => provider.finalize())

// verify with Pact, and reset expectations
afterEach(() => provider.verify())