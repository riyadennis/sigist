import logo from './logo.svg';
import './App.css';
import { useQuery, gql } from '@apollo/client';

const GET_LOCATIONS = gql`
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
`;

function DisplayLocations() {
  const { loading, error, data } = useQuery(GET_LOCATIONS);

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error : {error.message}</p>;

  return data.GetUser.map(({ id, firstName, lastName, email }) => (
    <body>
      <div class="App" key={id}>
        <h3>{firstName}</h3>
        <h3>{lastName}</h3>
        <h3 class="App-link">{email}</h3>
        <br />
    </div>
    </body>
    
  ));
}

export default function App() {
  return (
    <div >
      <h2>
        <div class="App-header">List of Users
        <span role="img" aria-label="rocket">
          🚀
          </span>
        </div>
      </h2>
      <br />
      <DisplayLocations />
    </div>
  );
}
