'use client'
import { Auth0Provider, useAuth0 } from '@auth0/auth0-react';

export default function Auth0({ children }) {
    return <div>
        <Auth0Provider
            domain="dev-fteqbjgrbz4fpbco.us.auth0.com"
            clientId='sLeeqUOaLtbW87QBFN8SJ8Q4WgXAPpuf'
            authorizationParams={{
                redirect_uri: 'http://localhost:3001',
                audience: "https://dev-fteqbjgrbz4fpbco.us.auth0.com/api/v2/",
                scope: "read:current_user update:current_user_metadata read:messages"
            }}
        >
            {children}
        </Auth0Provider>
    </div>
}