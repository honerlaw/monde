import * as React from "react";

export interface IPageNavbarProps {
    authPayload: boolean;
}

export class PageNavbar extends React.Component<IPageNavbarProps, {}> {

    public render(): JSX.Element {
        return <nav className="navbar navbar-expand-lg navbar-light">
            <a className="navbar-brand" href="#">package</a>
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
            <li className="nav-item">
                <a className="nav-link" href="/upload">upload</a>
            </li>
            <li className="nav-item">
                <a className="nav-link" href="/logout">logout</a>
            </li>
        </ul>;
    }

    private renderLoggedOut(): JSX.Element {
        return <ul className="navbar-nav ml-auto">
            <li className="nav-item">
                <a className="nav-link" href="/login">login</a>
            </li>
            <li className="nav-item">
                <a className="nav-link" href="/register">register</a>
            </li>
        </ul>;
    }

}
