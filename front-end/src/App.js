import { useState } from 'react';
import './App.css';
import { useQuery,useMutation, gql } from '@apollo/client';

const GET_USERS = gql`
  {
	GetUser(filter: {}){
		id
		firstName
		lastName
		email
	}
}
`;

const CREATE_USER = gql`
  mutation saveUser($input: CreateUserInput!){
	saveUser(input: $input){
		id
		email
		firstName
		lastName
		createAt
	}
}
`;

function DisplayUsers() {
  const { loading, error, data, refetch } = useQuery(GET_USERS);

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error : {error.message}</p>;
  if (data.GetUser.length === 0)  return <p>No Users</p>; 

  return data.GetUser.map(({ id, firstName, lastName, email }) => (
      <div class="App" key={id}>
        <h3>{firstName}</h3>
        <h3>{lastName}</h3>
        <h3 class="App-link">{email}</h3>
      <br />
    </div>
    
  ));
}


function DisplayForm() {  
  // Create User States
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [email, setEmail] = useState("");
  const [jobTitle, setJobTitle] = useState("");

  const [createUser] = useMutation(CREATE_USER, {
    variables: {
        input: {
          firstName: firstName,
          lastName: lastName,
          email: email,
          jobTitle: jobTitle,
        }
      }
  });
  

  return (
    <div class="container">
        <h1>Contact Form</h1>
      <hr/>
        <div>
        <label for="fristName"><b>First Name</b></label>
          <input
            name='fristName'
            placeholder='Given Name'
            onChange={(event) => {
              setFirstName(event.target.value);
            }}
          />
        </div>
        <div>
        <label for="lastName"><b>Last Name</b></label>
        <input
            name='lastName'
            placeholder='Surname'
            onChange={(event) => {
              setLastName(event.target.value);
            }}
          />
        </div>
        <div>
        <label for="email"><b>Email</b></label>
        <input
            name='email'
            placeholder='Email'
            onChange={(event) => {
              setEmail(event.target.value);
            }}
          />
        </div>
        <div>
        <label for="jobTitle"><b>Job Title</b></label>
        <input
            name='jobTitle'
            placeholder='Job Title'
            onChange={(event) => {
              setJobTitle(event.target.value);
            }}
          />
        </div>
        <button
          onClick={() => {
            createUser();
          }}
          class="registerbtn" >
          Create User
        </button>
            
    </div>
  );
}

export default function App() {
  return (
    <div>
      <DisplayForm />
      <h2>
        <div class="App-header">List of Users
        <span role="img" aria-label="rocket">
          🚀
          </span>
        </div>
      </h2>
      <DisplayUsers />
      
      <br />
    </div>
  );
}
