import * as React from "react";
import {registerComponent} from "preact-rpc";
import {Page} from "../page";

interface IProps {
    error: string;
}

export class LoginPage extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        return <Page>
            { this.props.error ? <span color="red">{this.props.error}</span> : null }
            <form method={"POST"} action={"/login"} id={"login-page"}>
                <input type={"text"} name={"username"} placeholder={"username"}/>
                <input type={"password"} name={"password"} placeholder={"password"}/>
                <button type={"submit"}>submit</button>
            </form>
        </Page>;
    }

}

registerComponent('LoginPage', LoginPage);
