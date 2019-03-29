import "babel-polyfill";
import "../../scss/global.scss"
import * as React from "react";
import {renderToString} from "react-dom/server";
import {Route, StaticRouter, Switch} from "react-router";
import {NotFoundPage} from "./components/pages/not-found-page";
import {LoginPage} from "./components/pages/login-page";
import {RegisterPage} from "./components/pages/register-page";
import {HomePage} from "./components/pages/home-page";
import {UploadListPage} from "./components/pages/upload-list-page";
import {Page} from "./components/page";
import {BrowserRouter} from "react-router-dom";
import {hydrate} from "react-dom";

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

global.main = function (opts: string, cbk: string) {
    const callback: (result: IResult) => void = global[cbk];
    const result: IResult = {
        error: null,
        html: null
    };
    try {
        const options: IOptions = JSON.parse(opts);
        try {
            result.html = renderToString(<StaticRouter location={options.url}>
                <Page options={opts} uploadForm={options.props.uploadForm} authPayload={options.props.authPayload}>
                    <Switch>
                        <Route exact path={"/"} render={() => <HomePage {...options.props}/>}/>
                        <Route path={"/user/login"} render={() => <LoginPage {...options.props}/>}/>
                        <Route path={"/user/register"} render={() => <RegisterPage {...options.props}/>}/>
                        <Route path={"/media/list"} render={() => <UploadListPage {...options.props}/>}/>
                        <Route render={() => <NotFoundPage {...options.props} />}/>
                    </Switch>
                </Page>
            </StaticRouter>)
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
    hydrate(<StaticRouter location={options.url}>
        <Page options={null} uploadForm={options.props.uploadForm} authPayload={options.props.authPayload}>
            <Switch>
                <Route exact path={"/"} render={() => <HomePage {...options.props}/>}/>
                <Route path={"/user/login"} render={() => <LoginPage {...options.props}/>}/>
                <Route path={"/user/register"} render={() => <RegisterPage {...options.props}/>}/>
                <Route path={"/media/list"} render={() => <UploadListPage {...options.props}/>}/>
                <Route render={() => <NotFoundPage {...options.props} />}/>
            </Switch>
        </Page>
    </StaticRouter>, document as any)
}