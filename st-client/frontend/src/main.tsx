import React from 'react'
import {createRoot} from 'react-dom/client'
import './style.css'
import App from './App'
import {createBrowserRouter, RouterProvider} from "react-router-dom";
import {Dashboard} from "./pages/dashboard";
import {Connections} from "./pages/connections";
import {Settings} from "./pages/settings";
import {Provider} from "react-redux";
import {store} from "./store/store";

const container = document.getElementById('root')

const root = createRoot(container!)

const router = createBrowserRouter([
    {
        path: "/",
        element: <App/>,
        children: [
            {
                // path: "dashboard",
                path: "/",
                element: <Dashboard/>
            },
            {
                path: "settings",
                element: <Settings/>
            },
            // {
            //     path: "connections",
            //     element: <Connections/>
            // }
        ]
    },
]);

root.render(
    <React.StrictMode>
        <Provider store={store}>
            <RouterProvider router={router} />
        </Provider>
    </React.StrictMode>
)
