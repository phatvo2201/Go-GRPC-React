import { useNavigate, Link } from "react-router-dom";
import { useCookies } from 'react-cookie';


const Home = () => {
    const [ removeCookie] = useCookies();

    const navigate = useNavigate();

    const logout = async () => {
        // if used in more components, this should be in context 
        // axios to /logout endpoint 
        navigate('/linkpage');
        removeCookie("token")

    }

    return (
        <section>
            <h1>Home</h1>
            <br />
            <p>You are logged in!</p>
            <br />
            <br />
            <br />
            <Link to="/wallet">Go to the Wallet</Link>
            <br />
            <Link to="/linkpage">Go to the link page</Link>
            <div className="flexGrow">
                <button onClick={logout}>Sign Out</button>
            </div>
        </section>
    )
}

export default Home
