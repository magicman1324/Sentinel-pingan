package com.pingan.monitor.deserializer;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.pingan.monitor.model.MetricEvent;
import org.apache.flink.api.common.serialization.DeserializationSchema;
import org.apache.flink.api.common.typeinfo.TypeInformation;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;

public class MetricDeserializer implements DeserializationSchema<MetricEvent> {

    private static final Logger LOG = LoggerFactory.getLogger(MetricDeserializer.class);
    private static final ObjectMapper mapper = new ObjectMapper();

    @Override
    public MetricEvent deserialize(byte[] message) {
        try {
            return mapper.readValue(message, MetricEvent.class);
        } catch (IOException e) {
            LOG.warn("dead-letter: invalid JSON metric, {} bytes, error={}", message.length, e.getMessage());
            return null; // Flink drops null records silently — dead-letter handled via metrics
        }
    }

    @Override
    public boolean isEndOfStream(MetricEvent nextElement) {
        return false;
    }

    @Override
    public TypeInformation<MetricEvent> getProducedType() {
        return TypeInformation.of(MetricEvent.class);
    }
}
