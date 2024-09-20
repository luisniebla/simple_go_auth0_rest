'use client'
import { useAuth0 } from "@auth0/auth0-react";
import { useState } from 'react'

const Login = () => {
  const { loginWithRedirect } = useAuth0();

  return (
    <button className='border-white border-2 p-1' onClick={() => loginWithRedirect()}>
      Login
    </button>
  )
}

const privateApiRequest = async (token) => {
  return fetch('http://localhost:8080/api/private', {
    headers: {
      'Authorization': `Bearer ${token}`,
      'Access-Control-Allow-Origin': '*'
    }
  }).then((resp) => {
    return resp.json()
  }).then((jsonResp => {
    return jsonResp
  }))
}


const AuthedComponent = ({ }) => {
  const { getAccessTokenSilently } = useAuth0()
  const [message, setMessage] = useState()

  return <><div>{message}</div><button onClick={async () => {
    const token = await getAccessTokenSilently()
    const json = await privateApiRequest(token)
    setMessage(json.message)
  }}>Request</button>
  </>

}

export default function Home() {
  const { user, logout } = useAuth0();

  return (
    <div className="grid items-center justify-items-center min-h-screen">
      <main className="">
        <h2>{user?.name}</h2>
        <p>{user?.email}</p>
        {user && <AuthedComponent />}
        {user && <button onClick={() => logout({ logoutParams: { returnTo: window.location.origin } })}>
          Log Out
        </button>}
        {!user && (
          <Login />
        )}
      </main>
    </div>
  )
}
