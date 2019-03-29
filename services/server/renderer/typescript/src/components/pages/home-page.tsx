import * as React from "react";

export class HomePage extends React.Component<{}, {}> {

    public render(): JSX.Element {
        return <div id={"home-page"}>
            <div className={"row"}>
                <div className={"col-sm-4 offset-sm-4"}>
                    Hello!
                </div>
            </div>
        </div>;
    }

}
