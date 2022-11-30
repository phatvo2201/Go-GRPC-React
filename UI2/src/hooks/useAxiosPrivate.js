import { axiosPrivate } from "../api/axios";
import { useEffect } from "react";
import useRefreshToken from "./useRefreshToken";
import { useCookies } from 'react-cookie';


const useAxiosPrivate = () => {
    const [cookies, setCookie] = useCookies(['user','token','rftoken','roles']);

    const refresh = useRefreshToken();
    let token =''
    if (!cookies.token){
        token  = ""
    }else{
        token = cookies.token
    }

    useEffect(() => {


        //use interceptor to add token
        const requestIntercept = axiosPrivate.interceptors.request.use(
            config => {
                if (!config.headers['Authorization']) {
                    config.headers['Authorization'] = `Bearer ${token}`;
                }
                console.log("request interceptor")
                return config;
            }, (error) => Promise.reject(error)
        );
        //check if exp token and add new token
        const responseIntercept = axiosPrivate.interceptors.response.use(
            response => response,
            async (error) => {
                const prevRequest = error?.config;
                if (error.response.status === 401 && !prevRequest?.sent) {
                    prevRequest.sent = true;
                    console.log("response interceptor")

                    const newAccessToken = await refresh();
                    prevRequest.headers['Authorization'] = `Bearer ${newAccessToken}`;

                    return axiosPrivate(prevRequest);
                }
                return Promise.reject(error);
            }
        );

        return () => {
            axiosPrivate.interceptors.request.eject(requestIntercept);
            axiosPrivate.interceptors.response.eject(responseIntercept);
        }
    }, [])

    return axiosPrivate;
}

export default useAxiosPrivate;