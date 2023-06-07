import { Matchers, GraphQLInteraction } from '@pact-foundation/pact'
import { provider } from './provider'
import graphql from 'graphQLAPI'

const { PACT_BROKER_TOKEN, PACT_BROKER_URL } = process.env;

const { like } = Matchers

const user = {
  id: like(123).contents,
  firstName: like('John').contents,
  lastName: like('Doe').contents,
  email: like('john@gmail.com').contents
}

describe('GraphQL', () => {
  describe('query users list', () => {
    beforeEach(() => {
      const graphqlQuery = new GraphQLInteraction()
        .uponReceiving('a list of users')
        .withRequest({
          path: '/graphql',
          method: 'POST'
        })
        .withOperation('GetUsers')
          .withQuery(`
            {
                GetUser(filter: {
                    firstName: "Aleena"
                }){
                    id
                    firstName
                    lastName
                    email
                }
            }
        `)
        .withVariables({})
        .willRespondWith({
          status: 200,
          headers: {
            'Content-Type': 'application/json; charset=utf-8'
          },
          body: {
            data: {
                  items: [
                    {  user  }
                  ]
            }
          }
        })
      return provider.addInteraction(graphqlQuery)
    })

    it('returns the correct response', async () => {
      expect(await graphql.GetUsers()).toEqual([user])
    })
  })
})
