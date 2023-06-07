import { render, screen } from '@testing-library/react';
import App from './App';
import { ApolloClient, InMemoryCache, ApolloProvider } from '@apollo/client';

test('renders learn react link', () => {
  const client = new ApolloClient({
    uri: 'http://localhost:8080/graphql',
    cache: new InMemoryCache(),
    Headers: {
      'Content-Type': 'application/json',
      "Access-Control-Allow-Origin": "http://localhost:8080/graphql",
    }
  });
  render( <ApolloProvider client={client}>
    <App />
  </ApolloProvider>,);
  const linkElement = screen.getByText(/users/i);
  expect(linkElement).toBeInTheDocument();
});
