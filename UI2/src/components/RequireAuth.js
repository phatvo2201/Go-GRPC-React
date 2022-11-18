import { useLocation, Navigate, Outlet } from "react-router-dom";
import useAuth from "../hooks/useAuth";
import { useCookies } from 'react-cookie';


const RequireAuth = ({ allowedRoles }) => {
    const { auth } = useAuth();
    const location = useLocation();
    const [cookies ] = useCookies(['user','token','roles']);

    const roles = cookies.roles






    return (
        roles?.roles?.find(role => allowedRoles?.includes(role))
        // allowedRoles.includes(roles)
            ? <Outlet />
            : auth?.user
                ? <Navigate to="/unauthorized" state={{ from: location }} replace />
                : <Navigate to="/login" state={{ from: location }} replace />
    );
}

export default RequireAuth;