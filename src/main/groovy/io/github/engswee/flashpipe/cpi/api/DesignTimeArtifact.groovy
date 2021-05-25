package io.github.engswee.flashpipe.cpi.api

import groovy.json.JsonBuilder
import groovy.json.JsonSlurper
import io.github.engswee.flashpipe.http.HTTPExecuter
import io.github.engswee.flashpipe.http.HTTPExecuterException
import org.slf4j.Logger
import org.slf4j.LoggerFactory

class DesignTimeArtifact {

    final HTTPExecuter httpExecuter

    static Logger logger = LoggerFactory.getLogger(DesignTimeArtifact)

    DesignTimeArtifact(HTTPExecuter httpExecuter) {
        this.httpExecuter = httpExecuter
    }

    String getVersion(String iFlowId, String iFlowVersion, boolean skipNotFoundException) {
        logger.info('Get Design time artifact')
        this.httpExecuter.executeRequest("/api/v1/IntegrationDesigntimeArtifacts(Id='$iFlowId',Version='$iFlowVersion')", ['Accept': 'application/json'])

        def code = this.httpExecuter.getResponseCode()
        logger.info("HTTP Response code = ${code}")
        if (code == 200) {
            def root = new JsonSlurper().parse(this.httpExecuter.getResponseBody())
            return root.d.Version
        } else if (skipNotFoundException && code == 404) {
            def root = new JsonSlurper().parse(this.httpExecuter.getResponseBody())
            if (root.error.message.value == 'Integration design time artifact not found') {
                return null
            } else {
                logger.info("Response body = ${this.httpExecuter.getResponseBody().getText('UTF8')}")
                throw new HTTPExecuterException("Get design time artifact call failed with response code = ${code}")
            }
        } else
            throw new HTTPExecuterException("Get design time artifact call failed with response code = ${code}")
    }

    byte[] download(String iFlowId, String iFlowVersion) {
        logger.info('Download Design time artifact')
        this.httpExecuter.executeRequest("/api/v1/IntegrationDesigntimeArtifacts(Id='$iFlowId',Version='$iFlowVersion')/\$value")

        def code = this.httpExecuter.getResponseCode()
        logger.info("HTTP Response code = ${code}")
        if (code == 200) {
            byte[] responseBody = this.httpExecuter.getResponseBody().getBytes()
            return responseBody
        } else
            throw new HTTPExecuterException("Download design time artifact call failed with response code = ${code}")
    }

    void update(String iFlowContent, String iFlowId, String iFlowName, String packageId, CSRFToken csrfToken) {
        // 1 - Get CSRF token
        String token = csrfToken ? csrfToken.get() : ''

        // 2 - Update IFlow
        updateArtifact(iFlowName, iFlowId, packageId, iFlowContent, token)
    }

    void delete(String iFlowId, CSRFToken csrfToken) {
        // 1 - Get CSRF token
        String token = csrfToken ? csrfToken.get() : ''

        // 2 - Update IFlow
        deleteArtifact(iFlowId, token)
    }

    String upload(String iFlowContent, String iFlowId, String iFlowName, String packageId, CSRFToken csrfToken) {
        // 1 - Get CSRF token
        String token = csrfToken ? csrfToken.get() : ''

        // 3 - Upload IFlow
        return uploadArtifact(iFlowName, iFlowId, packageId, iFlowContent, token)
    }

    void deploy(String iFlowId, CSRFToken csrfToken) {
        // 1 - Get CSRF token
        String token = csrfToken ? csrfToken.get() : ''

        // 2 - Deploy IFlow
        logger.info('Deploy design time artifact')
        this.httpExecuter.executeRequest('POST', '/api/v1/DeployIntegrationDesigntimeArtifact', ['x-csrf-token': token, 'Accept': 'application/json'], ['Id': "'${iFlowId}'", 'Version': "'active'"])
        def code = this.httpExecuter.getResponseCode()
        logger.info("HTTP Response code = ${code}")
        if (code != 202) {
            logger.info("Response body = ${this.httpExecuter.getResponseBody().getText('UTF8')}")
            throw new HTTPExecuterException("Deploy design time artifact call failed with response code = ${code}")
        }
    }

    private String constructPayload(String iFlowName, String iFlowId, String packageId, String iFlowContent) {
        def builder = new JsonBuilder()
        builder {
            'Name' iFlowName
            'Id' iFlowId
            'PackageId' packageId
            'ArtifactContent' iFlowContent
        }
        return builder.toString()
    }

    private void updateArtifact(String iFlowName, String iFlowId, String packageId, String iFlowContent, String token) {
        logger.info('Update design time artifact')
        def payload = constructPayload(iFlowName, iFlowId, packageId, iFlowContent)
        logger.debug("Request body = ${payload}")
        this.httpExecuter.executeRequest('PUT', "/api/v1/IntegrationDesigntimeArtifacts(Id='${iFlowId}',Version='active')", ['x-csrf-token': token, 'Accept': 'application/json'], null, payload, 'UTF-8', 'application/json')
        def code = this.httpExecuter.getResponseCode()
        logger.info("HTTP Response code = ${code}")
        if (code != 200) {
            logger.info("Response body = ${this.httpExecuter.getResponseBody().getText('UTF8')}")
            throw new HTTPExecuterException("Update design time artifact call failed with response code = ${code}")
        }
    }

    private void deleteArtifact(String iFlowId, String token) {
        logger.info('Delete existing design time artifact')
        this.httpExecuter.executeRequest('DELETE', "/api/v1/IntegrationDesigntimeArtifacts(Id='$iFlowId',Version='active')", ['x-csrf-token': token], null)
        def code = this.httpExecuter.getResponseCode()
        logger.info("HTTP Response code = ${code}")
        if (code != 200)
            throw new HTTPExecuterException("Delete design time artifact call failed with response code = ${code}")
    }

    private String uploadArtifact(String iFlowName, String iFlowId, String packageId, String iFlowContent, String token) {
        logger.info('Upload design time artifact')
        def payload = constructPayload(iFlowName, iFlowId, packageId, iFlowContent)
        logger.debug("Request body = ${payload}")

        this.httpExecuter.executeRequest('POST', '/api/v1/IntegrationDesigntimeArtifacts', ['x-csrf-token': token, 'Accept': 'application/json'], null, payload, 'UTF-8', 'application/json')
        def code = this.httpExecuter.getResponseCode()
        logger.info("HTTP Response code = ${code}")
        if (code != 201) {
            logger.info("Response body = ${this.httpExecuter.getResponseBody().getText('UTF8')}")
            throw new HTTPExecuterException("Upload design time artifact call failed with response code = ${code}")
        }

        return this.httpExecuter.getResponseBody().getText('UTF-8')
    }
}