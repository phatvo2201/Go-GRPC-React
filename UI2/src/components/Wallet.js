import { Link } from "react-router-dom"
import {  useState, useEffect } from "react";
import useAxiosPrivate from "../hooks/useAxiosPrivate";
import useRefreshToken  from "../hooks/useRefreshToken";
import { useCookies } from 'react-cookie';
import { axiosPrivate } from "../api/axios";




import { useLocation, useNavigate,  } from "react-router-dom";
import axios from '../api/axios';

 

import useAuth from "../hooks/useAuth";


const Wallet = () => {
    const [cookies, setCookie] = useCookies(['user','token','rftoken','roles']);

    const axiosPrivate = useAxiosPrivate();
    const refresh = useRefreshToken();

    const [wallet, setWallet] = useState();

    const navigate = useNavigate();
    const location = useLocation();

    const email = cookies.user.user

    useEffect(() => {
        let isMounted = true;
        const controller = new AbortController();
        const getUsers = async () => {
            try {
                const response = await axiosPrivate.post('http://localhost:8080/api/v1/userinfo/get_wallet',
                    JSON.stringify({ "Gmail":email}),
                    {
                        signal: controller.signal
                    }
                );
                isMounted && setWallet(response.data.balance);
            } catch (err) {
                console.error(err);
                navigate('/login', { state: { from: location }, replace: true });
            }
        }

        getUsers();

        return () => {
            isMounted = false;
            controller.abort();
        }
    }, [])

    return (
        <section>
            <h1>{email}</h1>


            <br />
            <p>This is your balance in wallet:</p>
            <h1>{wallet}</h1>

            <div className="flexGrow">
                <Link to="/">Home</Link>
            </div>
            <button onClick={()=>refresh()}>refresh</button>
        </section>
    )
}

export default Wallet
