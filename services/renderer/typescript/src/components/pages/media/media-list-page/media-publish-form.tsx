import * as React from "react";

interface IProps {
    videoId: string;
    canPublish: boolean;
    isPublished: boolean;
}

export class MediaPublishForm extends React.Component<IProps, {}> {

    public render(): JSX.Element | null {
        if (!this.props.canPublish) {
            return null;
        }

        return <form method={"POST"} action={"/media/publish"} className={"publish-form"}>
            <input type={"hidden"} name={"video_id"} value={this.props.videoId}/>
            <button className="btn btn-primary" type={"submit"}>
                {this.props.isPublished ? "unpublish" : "publish"}
            </button>
        </form>
    }

}