import { useLocation, Navigate, Outlet } from "react-router-dom";
import useAuth from "../hook/useAuth";
import Spinner from "./spinner/Spinner";


//If the user is not aiuthenticated and tries to access a protected route, we redirect him to login route
const RequiredAuth = () => {
    const {auth, loading} = useAuth();
    const location = useLocation();

    if (loading) {
        return (<Spinner/>)
    }

    return auth ? (
        <Outlet/>
    ) : (
        <Navigate to = '/login' state={{from: location}} replace />
    );
};

export default RequiredAuth;