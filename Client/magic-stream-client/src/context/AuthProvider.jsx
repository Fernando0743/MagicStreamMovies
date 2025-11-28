import { Children, createContext, useEffect, useState } from "react";

const AuthContext = createContext({});

//Children is a prop that represents all the components in our application
export const AuthProvider = ({children}) => {
    const [auth, setAuth] = useState();
    const [loading, setLoading] = useState(true);

    //Once componente starts loading, read local storage to see if user is authenticated
    useEffect(() => {
        try{
            const storedUser = localStorage.getItem('user');
            if (storedUser) {
                const parsedUser = JSON.parse(storedUser);
                setAuth(parsedUser)
            }
        } catch(error) {
            console.error("Failed to parse user from local storage: ", error)
        } finally {
            setLoading(false);
        }
    }, []);

    //This useEffect runs only when auth variable changes, so auth will change when user logs in or logs out
    useEffect(() => {
        if(auth){
            localStorage.setItem('user', JSON.stringify(auth));
        } else {
            localStorage.removeItem('user');
        }
    },[auth])


    return (
        /*Le tenemos que pasar setAuth y Auth para darle a entender que esas dos cosas estan disponibles para todos nuestros children de la aplicacion */
        <AuthContext.Provider value = {{auth,setAuth, loading}}>
            {children}
        </AuthContext.Provider>
    )
}

export default AuthContext;