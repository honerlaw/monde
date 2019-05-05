import * as React from "react";
import {IGlobalProps} from "../global-props";
import {UploadForm} from "./media/upload-form";

export interface IPageNavbarProps extends IGlobalProps {
    uploadForm: any;
}

/**
 * @TODO
 * - fix dropdown stuff to actually work... since when the window is small the dropdown is broken
 */
export class PageNavbar extends React.Component<IPageNavbarProps, {}> {

    public render(): JSX.Element {
        return <nav className="navbar navbar-expand-lg navbar-light">
            <a className="navbar-brand" href="/">vueon</a>
            <button className="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarNavDropdown"
                    aria-controls="navbarNavDropdown" aria-expanded="false" aria-label="Toggle navigation">
                <span className="navbar-toggler-icon" />
            </button>
            <div className="collapse navbar-collapse" id="navbarNavDropdown">
                { this.props.authPayload ? this.renderLoggedIn() : this.renderLoggedOut()}
            </div>
        </nav>;
    }

    private renderLoggedIn(): JSX.Element {
        return <ul className="navbar-nav ml-auto">
            <li className={"nav-item"}>
                <UploadForm form={this.props.uploadForm}/>
            </li>
            <li className="nav-item">
                <a className="nav-link" href="/media/list">my videos</a>
            </li>
            <li className="nav-item">
                <a className="nav-link" href="/user">account</a>
            </li>
            <li className="nav-item">
                <a className="nav-link" href="/user/logout">logout</a>
            </li>
        </ul>;
    }

    private renderLoggedOut(): JSX.Element {
        return <ul className="navbar-nav ml-auto">
            <li className="nav-item">
                <a className="nav-link" href="/user/login">login</a>
            </li>
            <li className="nav-item">
                <a className="nav-link" href="/user/register">register</a>
            </li>
        </ul>;
    }

}
