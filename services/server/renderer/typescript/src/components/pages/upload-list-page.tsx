import * as React from "react";
import {IUploadForm} from "./upload-list-page/upload-form";
import {PendingUploadItem} from "./upload-list-page/pending-upload-item";
import {UploadItem} from "./upload-list-page/upload-item";
import {IGlobalProps} from "../../global-props";

export interface IUploadInfo {
    videoId: string;
    canPublish: boolean;
    info: {
        title: string;
        description: string;
        status: string;
        hashtags: string[];
        published: boolean;
    };
    thumbs: string[];
    videos: Array<{
        type: string;
        url: string;
    }>
}

interface IProps extends IGlobalProps {
    uploads: IUploadInfo[];
    uploadForm: IUploadForm;
}

/**
 * @todo
 * - display thumbnail after upload so the user can see what it is
 */
export class UploadListPage extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        return <div id={"upload-list-page"}>
            <ol className={"upload-list"}>
                {this.props.uploads.map((upload: IUploadInfo): JSX.Element => {
                    if (upload.info.status !== "Complete") {
                        return <PendingUploadItem key={upload.videoId} status={upload.info.status}/>;
                    }
                    return <UploadItem key={upload.videoId} upload={upload}/>;
                })}
            </ol>
        </div>;
    }

}
