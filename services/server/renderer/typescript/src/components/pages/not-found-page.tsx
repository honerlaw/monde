import * as React from "react";

export class NotFoundPage extends React.Component<{}, {}> {

    public render(): JSX.Element {
        return <div id={"not-found-page"}>
            <span>Page Not Found!</span>
        </div>;
    }

}
