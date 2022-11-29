import { useLocation, Navigate, Outlet } from "react-router-dom";
import { useCookies } from 'react-cookie';


const RequireAuth = ({ allowedRoles }) => {
    const location = useLocation();
    const [cookies] = useCookies(['roles', 'token'])


    const roles = cookies.roles

    return (
        roles?.roles?.find(role => allowedRoles?.includes(role))
            ? <Outlet />
            : cookies?.token
                ? <Navigate to="/unauthorized" state={{ from: location }} replace />
                : <Navigate to="/login" state={{ from: location }} replace />
    );
}

export default RequireAuth;