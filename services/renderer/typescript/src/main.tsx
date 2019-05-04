import "babel-polyfill";
import "./main.scss"
import * as React from "react";
import {renderToString} from "react-dom/server";
import {Route, RouteComponentProps, StaticRouter, Switch} from "react-router";
import {NotFoundPage} from "./components/pages/error/not-found-page";
import {LoginPage} from "./components/pages/login-page";
import {RegisterPage} from "./components/pages/register-page";
import {HomePage} from "./components/pages/home-page";
import {Page} from "./components/page";
import {hydrate} from "react-dom";
import {MediaListPage} from "./components/pages/media/media-list-page";
import {ComponentClass} from "react";
import {MediaViewPage} from "./components/pages/media/media-view-page";
import {UserPage} from "./components/pages/user-page";

declare const global: any;

interface IResult {
    error: string | null;
    html: string | null;
}

interface IOptions {
    url: string;
    headers: { [key: string]: string };
    props: any;
}

function renderRoute(Component: any, props: any): (props: RouteComponentProps<any>) => JSX.Element {
    return (routerProps: RouteComponentProps<any>): JSX.Element => {
        const allProps = {...routerProps, ...props};
        return <Component {...allProps} />;
    }
}

/**
 * Render all routes in a static router
 *
 * @param {IOptions} options
 * @param {string|null} opts
 */
function element(options: IOptions, opts: string | null): JSX.Element {
    return <StaticRouter location={options.url}>
        <Page options={opts} uploadForm={options.props.uploadForm} authPayload={options.props.authPayload}>
            <Switch>
                <Route exact path={"/"} render={renderRoute(HomePage, options.props)}/>
                <Route exact path={"/user"} render={renderRoute(UserPage, options.props)}/>
                <Route path={"/user/login"} render={renderRoute(LoginPage, options.props)}/>
                <Route path={"/user/register"} render={renderRoute(RegisterPage, options.props)}/>
                <Route path={"/media/list"} render={renderRoute(MediaListPage, options.props)}/>
                <Route path={"/media/view/:mediaId"} render={renderRoute(MediaViewPage, options.props)}/>
                <Route render={renderRoute(NotFoundPage, options.props)}/>
            </Switch>
        </Page>
    </StaticRouter>;
}

/**
 * Called by goja to handle server side rendering
 *
 * @param {IOptions} opts - options passed from the go server, its a string because of weird issues with goja's ToValue
 * @param {string} cbk - name of the callback in the global scope
 */
global.main = function (opts: string, cbk: string) {
    const callback: (result: IResult) => void = global[cbk];
    const result: IResult = {
        error: null,
        html: null
    };
    try {
        const options: IOptions = JSON.parse(opts);
        try {
            result.html = renderToString(element(options, opts))
        } catch (err) {
            result.error = err.toString();
        }
    } catch (err) {
        result.error = err.toString()
    }

    return callback(result);
};

// when run in the browser, this will be executed and hydrate the html with the event listeners / etc
if (global.window !== undefined) {
    const options: any = (window as any).hydrateOptions;
    hydrate(element(options, null), document as any)
}
