// Libraries
import React, {PureComponent} from 'react'
import {connect} from 'react-redux'
import _ from 'lodash'

// Components
import {Form} from '@influxdata/clockface'
import {ErrorHandling} from 'src/shared/decorators/errors'
import StreamingSelector from 'src/dataLoaders/components/collectorsWizard/select/StreamingSelector'
import OnboardingButtons from 'src/onboarding/components/OnboardingButtons'
import FancyScrollbar from 'src/shared/components/fancy_scrollbar/FancyScrollbar'

// Actions
import {
  addPluginBundleWithPlugins,
  removePluginBundleWithPlugins,
} from 'src/dataLoaders/actions/dataLoaders'
import {setBucketInfo} from 'src/dataLoaders/actions/steps'

// Types
import {Bucket} from '@influxdata/influx'
import {ComponentStatus} from '@influxdata/clockface'
import {CollectorsStepProps} from 'src/dataLoaders/components/collectorsWizard/CollectorsWizard'
import {TelegrafPlugin, BundleName} from 'src/types/dataLoaders'
import {AppState} from 'src/types'

export interface OwnProps extends CollectorsStepProps {
  buckets: Bucket[]
}

export interface StateProps {
  bucket: string
  telegrafPlugins: TelegrafPlugin[]
  pluginBundles: BundleName[]
}

export interface DispatchProps {
  onAddPluginBundle: typeof addPluginBundleWithPlugins
  onRemovePluginBundle: typeof removePluginBundleWithPlugins
  onSetBucketInfo: typeof setBucketInfo
}

type Props = OwnProps & StateProps & DispatchProps

@ErrorHandling
export class SelectCollectorsStep extends PureComponent<Props> {
  public render() {
    return (
      <Form
        onSubmit={this.props.onIncrementCurrentStepIndex}
        className="data-loading--form"
      >
        <FancyScrollbar
          autoHide={false}
          className="data-loading--scroll-content"
        >
          <div>
            <h3 className="wizard-step--title">What do you want to monitor?</h3>
            <h5 className="wizard-step--sub-title">
              Telegraf is a plugin-based data collection agent which writes
              metrics to a bucket in InfluxDB
            </h5>
          </div>
          {!!this.props.bucket && (
            <StreamingSelector
              pluginBundles={this.props.pluginBundles}
              telegrafPlugins={this.props.telegrafPlugins}
              onTogglePluginBundle={this.handleTogglePluginBundle}
              buckets={this.props.buckets}
              selectedBucketName={this.props.bucket}
              onSelectBucket={this.handleSelectBucket}
            />
          )}
        </FancyScrollbar>
        <OnboardingButtons
          autoFocusNext={true}
          nextButtonStatus={this.nextButtonStatus}
          className="data-loading--button-container"
        />
      </Form>
    )
  }

  private get nextButtonStatus(): ComponentStatus {
    const {telegrafPlugins, buckets} = this.props

    if (!buckets || !buckets.length) {
      return ComponentStatus.Disabled
    }

    if (!telegrafPlugins.length) {
      return ComponentStatus.Disabled
    }

    return ComponentStatus.Default
  }

  private handleSelectBucket = (bucket: Bucket) => {
    const {orgID, id, name} = bucket

    this.props.onSetBucketInfo(orgID, name, id)
  }

  private handleTogglePluginBundle = (
    bundle: BundleName,
    isSelected: boolean
  ) => {
    if (isSelected) {
      this.props.onRemovePluginBundle(bundle)

      return
    }

    this.props.onAddPluginBundle(bundle)
  }
}

const mstp = ({
  dataLoading: {
    dataLoaders: {telegrafPlugins, pluginBundles},
    steps: {bucket},
  },
}: AppState): StateProps => ({
  telegrafPlugins,
  bucket,
  pluginBundles,
})

const mdtp: DispatchProps = {
  onAddPluginBundle: addPluginBundleWithPlugins,
  onRemovePluginBundle: removePluginBundleWithPlugins,
  onSetBucketInfo: setBucketInfo,
}

export default connect<StateProps, DispatchProps, OwnProps>(
  mstp,
  mdtp
)(SelectCollectorsStep)
